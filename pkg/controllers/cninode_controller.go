/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/ipamd"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/ipamd/datastore"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/utils/logger"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/version"
	"github.com/aws/amazon-vpc-cni-k8s/utils/imds"
	"strings"

	cniNode "github.com/aws/amazon-vpc-cni-k8s/pkg/apis/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PolicyEndpointsReconciler reconciles a PolicyEndpoints object
type CNINodeReconciler struct {
	k8sClient client.Client
	scheme    *runtime.Scheme
	//Primary IP of EC2 instance
	nodeIP string
	//Node Name
	nodeName string
	//Map from ENI ID to primary IP of that ENI
	primaryIP map[string]string
	//Primary ENI MAC
	primaryMAC string
	//Primary ENI ID
	primaryENIID string
	//IPV4 VPC CIDRs
	v4VPCCIDRs []string
	//IPV6 VPC CIDRs
	v6VPCCIDRs []string
	//IPAM's internal datastore
	datastore *datastore.DataStore

	isNodeInitialized bool

	ipamContext *ipamd.IPAMContext

	//Logger
	log logger.Logger
}

// NewPolicyEndpointsReconciler constructs new PolicyEndpointReconciler
func NewCNINodeReconciler(k8sClient client.Client, log logger.Logger, enableIPv6 bool) (*CNINodeReconciler, error) {
	c := &CNINodeReconciler{
		k8sClient: k8sClient,
		log:       log,
	}

	c.nodeName, _ = imds.GetMetaData("hostname")
	c.primaryMAC, _ = imds.GetMetaData("mac")
	if !enableIPv6 {
		c.nodeIP, _ = imds.GetMetaData("local-ipv4")

	} else {
		c.nodeIP, _ = imds.GetMetaData("ipv6")
	}

	//Initialize Datastore
	checkPointer := datastore.NewJSONFile("/var/run/aws-node/ipam.json")
	c.datastore = datastore.NewDataStore(log, checkPointer, true)

	c.log.Infof("Initialized CNI Node reconciler for:- %s", c.nodeName)
	return c, nil
}

//+kubebuilder:rbac:groups=eks.amazonaws.com,resources=cninodes,verbs=get;list;watch
//+kubebuilder:rbac:groups=eks.amazonaws.com,resources=cninodes/status,verbs=get

func (c *CNINodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	c.log.Infof("Received a new reconcile request: %s", req)

	if err := c.reconcile(ctx, req); err != nil {
		c.log.Errorf("Reconcile error")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (c *CNINodeReconciler) reconcile(ctx context.Context, req ctrl.Request) error {
	cniNode := &cniNode.CNINode{}
	c.log.Infof("Get CNI Node object for: %s", req.NamespacedName.Name)
	//Check if this resource belongs to this node
	if req.NamespacedName.Name != c.nodeName {
		c.log.Infof("CNI Node resource doesn't belong to the node. Skipping....")
		return nil
	}

	if err := c.k8sClient.Get(ctx, req.NamespacedName, cniNode); err != nil {
		if apierrors.IsNotFound(err) {
			c.log.Errorf("Unable to get CNI Node object for: %s", req.NamespacedName.Name)
			return nil
		}
		c.log.Errorf("Unable to get CNI Node resource")
		return nil
	}
	/* DELETE should never happen for node specific CNI Node object until the node is deleted.
	if !cniNode.DeletionTimestamp.IsZero() {
		return c.cleanUpPolicyEndpoint(ctx, req)
	}
	*/

	c.log.Infof("Successfully derived CNI Node resource")
	return c.reconcileCNINode(ctx, cniNode)
}

func (c *CNINodeReconciler) reconcileCNINode(ctx context.Context, cniNode *cniNode.CNINode) error {
	var err error
	var networkInterfaces []ipamd.ENIMetadata
	//ipamContext := &ipamd.IPAMContext{}
	//Derive ENI, IPs/Prefixes and update the datastore
	if networkInterfaces, err = c.deriveNodeNetworkingInfo(ctx, cniNode); err != nil {
		c.log.Errorf("Unable to derive Node Networking Info")
		return nil
	}

	//Configure Primary ENI
	if !c.isNodeInitialized {
		if c.ipamContext, err = ipamd.New(c.k8sClient, true, false, c.nodeName, c.v4VPCCIDRs, c.v6VPCCIDRs,
			c.nodeIP, c.primaryMAC, c.primaryENIID, c.datastore, networkInterfaces); err != nil {
			c.log.Errorf("IPAMD Initialization failed")
		}
		go c.ipamContext.RunRPCHandler(version.Version)
		c.isNodeInitialized = true
	}
	return nil
}

func (c *CNINodeReconciler) deriveNodeNetworkingInfo(ctx context.Context, cniNode *cniNode.CNINode) ([]ipamd.ENIMetadata, error) {
	var err error
	var networkInterfaces []ipamd.ENIMetadata

	//Derive VPC CIDR info. Required for configuring relevant forwarding rules
	if err = c.deriveVPCCIDRs(ctx, cniNode); err != nil {
		c.log.Errorf("Unable to derive VPC CIDRs for the node")
		return nil, nil
	}

	//Derive Primary ENI IP and MAC
	if err = c.derivePrimaryENIInfo(ctx, cniNode); err != nil {
		c.log.Errorf("Unable to derive Primary ENI Info for the node")
		return nil, nil
	}

	//Derive information about ENIs and IPs assigned to them
	if networkInterfaces, err = c.deriveENIInfo(ctx, cniNode); err != nil {
		c.log.Errorf("Unable to derive IP CIDRs Info for the node")
		return nil, nil
	}

	return networkInterfaces, nil
}

func (c *CNINodeReconciler) deriveVPCCIDRs(ctx context.Context, cniNode *cniNode.CNINode) error {
	c.v4VPCCIDRs = cniNode.Status.VPCCIDRs
	c.v6VPCCIDRs = cniNode.Status.VPCCIDRs

	//Check if relevant CIDR is filled based on the IP Family of the cluster and error out if that's not the case

	return nil
}

func (c *CNINodeReconciler) derivePrimaryENIInfo(ctx context.Context, cniNode *cniNode.CNINode) error {
	c.primaryENIID = cniNode.Spec.PrimaryNetworkInterfaceID
	if len(strings.TrimSpace(c.primaryENIID)) == 0 {
		c.log.Errorf("Primary ENI Info is missing. CNI will not be able to operate...")
		return nil
	}

	c.log.Infof("Primary ENI ID: %s", c.primaryENIID)
	c.log.Infof("Primary IP: %s", c.nodeIP)
	return nil
}

func (c *CNINodeReconciler) deriveENIInfo(ctx context.Context, cniNode *cniNode.CNINode) ([]ipamd.ENIMetadata, error) {
	c.log.Infof("Derive CIDRs")
	var localNetworkInterfaces []ipamd.ENIMetadata
	for _, networkInterfaceInfo := range cniNode.Status.NetworkInterfaces {
		c.log.Infof("Deriving Network Interface info for ID: %s", networkInterfaceInfo.ID)
		var eniHostCIDRs, eniPrefixCIDRs, CoolDownHostCIDRs, CoolDownPrefixCIDRs []string
		for _, Cidr := range networkInterfaceInfo.CIDRs {
			c.log.Infof("IP CIDR: %s", string(Cidr))
			ipAddress, isPrefix := c.IsNonHostCIDR(string(Cidr))
			if isPrefix {
				eniPrefixCIDRs = append(eniPrefixCIDRs, string(Cidr))
			} else {
				eniHostCIDRs = append(eniHostCIDRs, ipAddress)
			}
		}

		for _, CoolDownCIDR := range networkInterfaceInfo.CoolDownCIDRs {
			for _, Cidr := range CoolDownCIDR.CIDRs {
				c.log.Infof("CoolDown CIDR: %s", string(Cidr))
				ipAddress, isPrefix := c.IsNonHostCIDR(string(Cidr))
				if isPrefix {
					CoolDownPrefixCIDRs = append(CoolDownPrefixCIDRs, string(Cidr))
				} else {
					CoolDownHostCIDRs = append(CoolDownHostCIDRs, ipAddress)
				}
			}
		}
		c.log.Infof("Total number of Host CIDRs attached: %d; Prefix CIDRs: %d", len(eniHostCIDRs), len(eniPrefixCIDRs))
		c.log.Infof("Total number of Host Cool Down CIDRs: %d; "+
			"Prefix CIDRs: %d", len(CoolDownHostCIDRs), len(CoolDownPrefixCIDRs))
		localNetworkInterfaces = append(localNetworkInterfaces,
			ipamd.ENIMetadata{
				ENIID:               networkInterfaceInfo.ID,
				MAC:                 c.primaryMAC,
				PrimaryIP:           string(networkInterfaceInfo.PrimaryCIDR),
				DeviceNumber:        networkInterfaceInfo.DeviceIndex,
				NetworkCardIndex:    networkInterfaceInfo.NetworkCardIndex,
				SubnetIPv4CIDR:      string(networkInterfaceInfo.SubnetCIDR),
				SubnetIPv6CIDR:      string(networkInterfaceInfo.SubnetCIDR),
				IPv4Addresses:       eniHostCIDRs,
				IPv4Prefixes:        eniPrefixCIDRs,
				CoolDownHostCIDRs:   CoolDownHostCIDRs,
				CoolDownPrefixCIDRs: CoolDownPrefixCIDRs,
			})
	}

	//Update Local data store
	if c.isNodeInitialized {
		c.log.Infof("Setup ENIs....")
		for _, eniMetadata := range localNetworkInterfaces {
			c.log.Infof("Configuring ENI: %s", eniMetadata.ENIID)
			err := c.ipamContext.SetupENI(eniMetadata, false, false, c.primaryENIID)
			if err != nil {
				c.log.Errorf("failed configuring ENI: %s", eniMetadata.ENIID)
			}
			//Purge the CIDRs in cool down list from the datastore
			err = c.ipamContext.VerifyAndDeletePrefixesFromDatastore(eniMetadata.ENIID, eniMetadata.CoolDownHostCIDRs,
				eniMetadata.CoolDownPrefixCIDRs)
		}
	}

	return localNetworkInterfaces, nil
}

func (c *CNINodeReconciler) IsNonHostCIDR(ipAddr string) (string, bool) {
	ipSplit := strings.Split(ipAddr, "/")
	//Ignore Catch All IP entry as well
	if ipSplit[1] != "32" && ipSplit[1] != "128" && ipSplit[1] != "0" {
		c.log.Infof("Non Host IP. Return")
		return ipSplit[0], true
	}
	return ipSplit[0], false
}

// SetupWithManager sets up the controller with the Manager.
func (c *CNINodeReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	c.log.Infof("Setting up the controller...")
	return ctrl.NewControllerManagedBy(mgr).
		For(&cniNode.CNINode{}).
		Complete(c)
}
