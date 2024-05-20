// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package ipamd

import (
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/aws/amazon-vpc-cni-k8s/pkg/ipamd/datastore"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/networkutils"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/utils/logger"
	"github.com/aws/amazon-vpc-cni-k8s/utils"
	"github.com/aws/amazon-vpc-cni-k8s/utils/prometheusmetrics"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

// The package ipamd is a long running daemon which manages a warm pool of available IP addresses.
// It also monitors the size of the pool, dynamically allocates more ENIs when the pool size goes below
// the minimum threshold and frees them back when the pool size goes above max threshold.

const (
	ipPoolMonitorInterval       = 5 * time.Second
	maxRetryCheckENI            = 5
	eniAttachTime               = 10 * time.Second
	nodeIPPoolReconcileInterval = 60 * time.Second
	decreaseIPPoolInterval      = 30 * time.Second

	// ipReconcileCooldown is the amount of time that an IP address must wait until it can be added to the data store
	// during reconciliation after being discovered on the EC2 instance metadata.
	ipReconcileCooldown = 60 * time.Second

	// This environment variable is used to specify the desired number of free IPs always available in the "warm pool".
	// When it is not set, ipamd defaults to use all available IPs per ENI for that instance type.
	// For example, for a m4.4xlarge node,
	//     If WARM_IP_TARGET is set to 1, and there are 9 pods running on the node, ipamd will try
	//     to make the "warm pool" have 10 IP addresses with 9 being assigned to pods and 1 free IP.
	//
	//     If "WARM_IP_TARGET is not set, it will default to 30 (which the maximum number of IPs per ENI).
	//     If there are 9 pods running on the node, ipamd will try to make the "warm pool" have 39 IPs with 9 being
	//     assigned to pods and 30 free IPs.
	envWarmIPTarget = "WARM_IP_TARGET"
	noWarmIPTarget  = 0

	// This environment variable is used to specify the desired minimum number of total IPs.
	// When it is not set, ipamd defaults to 0.
	// For example, for a m4.4xlarge node,
	//     If WARM_IP_TARGET is set to 1 and MINIMUM_IP_TARGET is set to 12, and there are 9 pods running on the node,
	//     ipamd will make the "warm pool" have 12 IP addresses with 9 being assigned to pods and 3 free IPs.
	//
	//     If "MINIMUM_IP_TARGET is not set, it will default to 0, which causes WARM_IP_TARGET settings to be the
	//	   only settings considered.
	envMinimumIPTarget = "MINIMUM_IP_TARGET"
	noMinimumIPTarget  = 0

	// This environment is used to specify the desired number of free ENIs along with all of its IP addresses
	// always available in "warm pool".
	// When it is not set, it is default to 1.
	//
	// when "WARM_IP_TARGET" is defined, ipamd will use behavior defined for "WARM_IP_TARGET".
	//
	// For example, for a m4.4xlarge node
	//     If WARM_ENI_TARGET is set to 2, and there are 9 pods running on the node, ipamd will try to
	//     make the "warm pool" to have 2 extra ENIs and its IP addresses, in other words, 90 IP addresses
	//     with 9 IPs assigned to pods and 81 free IPs.
	//
	//     If "WARM_ENI_TARGET" is not set, it defaults to 1, so if there are 9 pods running on the node,
	//     ipamd will try to make the "warm pool" have 1 extra ENI, in other words, 60 IPs with 9 already
	//     being assigned to pods and 51 free IPs.
	envWarmENITarget     = "WARM_ENI_TARGET"
	defaultWarmENITarget = 1

	// This environment variable is used to specify the maximum number of ENIs that will be allocated.
	// When it is not set or less than 1, the default is to use the maximum available for the instance type.
	//
	// The maximum number of ENIs is in any case limited to the amount allowed for the instance type.
	envMaxENI     = "MAX_ENI"
	defaultMaxENI = -1

	// This environment is used to specify whether Pods need to use a security group and subnet defined in an ENIConfig CRD.
	// When it is NOT set or set to false, ipamd will use primary interface security group and subnet for Pod network.
	envCustomNetworkCfg = "AWS_VPC_K8S_CNI_CUSTOM_NETWORK_CFG"

	// This environment variable specifies whether IPAMD should allocate or deallocate ENIs on a non-schedulable node (default false).
	envManageENIsNonSchedulable = "AWS_MANAGE_ENIS_NON_SCHEDULABLE"

	// This environment is used to specify whether we should use enhanced subnet selection or not when creating ENIs (default true).
	envSubnetDiscovery = "ENABLE_SUBNET_DISCOVERY"

	// eniNoManageTagKey is the tag that may be set on an ENI to indicate ipamd
	// should not manage it in any form.
	eniNoManageTagKey = "node.k8s.amazonaws.com/no_manage"

	// disableENIProvisioning is used to specify that ENIs do not need to be synced during initializing a pod.
	envDisableENIProvisioning = "DISABLE_NETWORK_RESOURCE_PROVISIONING"

	// disableLeakedENICleanup is used to specify that the task checking and cleaning up leaked ENIs should not be run.
	envDisableLeakedENICleanup = "DISABLE_LEAKED_ENI_CLEANUP"

	// Specify where ipam should persist its current IP<->container allocations.
	envBackingStorePath     = "AWS_VPC_K8S_CNI_BACKING_STORE"
	defaultBackingStorePath = "/var/run/aws-node/ipam.json"

	// envEnablePodENI is used to attach a Trunk ENI to every node. Required in order to give Branch ENIs to pods.
	envEnablePodENI = "ENABLE_POD_ENI"

	// envNodeName will be used to store Node name
	envNodeName = "MY_NODE_NAME"

	//envEnableIpv4PrefixDelegation is used to allocate /28 prefix instead of secondary IP for an ENI.
	envEnableIpv4PrefixDelegation = "ENABLE_PREFIX_DELEGATION"

	//envWarmPrefixTarget is used to keep a /28 prefix in warm pool.
	envWarmPrefixTarget     = "WARM_PREFIX_TARGET"
	defaultWarmPrefixTarget = 0

	//envEnableIPv4 - Env variable to enable/disable IPv4 mode
	envEnableIPv4 = "ENABLE_IPv4"

	//envEnableIPv6 - Env variable to enable/disable IPv6 mode
	envEnableIPv6 = "ENABLE_IPv6"

	ipV4AddrFamily = "4"
	ipV6AddrFamily = "6"

	// insufficientCidrErrorCooldown is the amount of time reconciler will wait before trying to fetch
	// more IPs/prefixes for an ENI. With InsufficientCidr we know the subnet doesn't have enough IPs so
	// instead of retrying every 5s which would lead to increase in EC2 AllocIPAddress calls, we wait for
	// 120 seconds for a retry.
	insufficientCidrErrorCooldown = 120 * time.Second

	// envManageUntaggedENI is used to determine if untagged ENIs should be managed or unmanaged
	envManageUntaggedENI = "MANAGE_UNTAGGED_ENI"

	eniNodeTagKey = "node.k8s.amazonaws.com/instance_id"

	// envAnnotatePodIP is used to annotate[vpc.amazonaws.com/pod-ips] pod's with IPs
	// Ref : https://github.com/projectcalico/calico/issues/3530
	// not present; in which case we fall back to the k8s podIP
	// Present and set to an IP; in which case we use it
	// Present and set to the empty string, which we use to mean "CNI DEL had occurred; networking has been removed from this pod"
	// The empty string one helps close a trace at pod shutdown where it looks like the pod still has its IP when the IP has been released
	envAnnotatePodIP = "ANNOTATE_POD_IP"

	// aws error codes for insufficient IP address scenario
	INSUFFICIENT_CIDR_BLOCKS    = "InsufficientCidrBlocks"
	INSUFFICIENT_FREE_IP_SUBNET = "InsufficientFreeAddressesInSubnet"

	// envEnableNetworkPolicy is used to enable IPAMD/CNI to send pod create events to network policy agent.
	envNetworkPolicyMode     = "NETWORK_POLICY_ENFORCING_MODE"
	defaultNetworkPolicyMode = "standard"
	hostMask                 = "/32"
)

var log = logger.Get()

var (
	prometheusRegistered = false
)

// IPAMContext contains node level control information
type IPAMContext struct {
	dataStore     *datastore.DataStore
	k8sClient     client.Client
	enableIPv4    bool
	enableIPv6    bool
	networkClient networkutils.NetworkAPIs

	myNodeName            string
	primaryIP             map[string]string // primaryIP is a map from ENI ID to primary IP of that ENI
	terminating           int32             // Flag to warn that the pod is about to shut down.
	enablePodIPAnnotation bool
	maxPods               int // maximum number of pods that can be scheduled on the node
	networkPolicyMode     string
	v4VPCCIDRs            []string
	v6VPCCIDRs            []string

	// reconcileCooldownCache keeps timestamps of the last time an IP address was unassigned from an ENI,
	// so that we don't reconcile and add it back too quickly if IMDS lags behind reality.
	reconcileCooldownCache ReconcileCooldownCache
}

// ENIMetadata contains information about an ENI
type ENIMetadata struct {
	// ENIID is the id of network interface
	ENIID string

	// MAC is the mac address of network interface
	MAC string

	// DeviceNumber is the  device number of network interface
	DeviceNumber int32 // 0 means it is primary interface

	// NetworkCardIndex
	NetworkCardIndex int32

	//Primary IP of the ENI
	PrimaryIP string

	// SubnetIPv4CIDR is the IPv4 CIDR of network interface
	SubnetIPv4CIDR string

	// SubnetIPv6CIDR is the IPv6 CIDR of network interface
	SubnetIPv6CIDR string

	// The ip addresses allocated for the network interface
	IPv4Addresses []string

	// IPv4 Prefixes allocated for the network interface
	IPv4Prefixes []string

	// IPv6 addresses allocated for the network interface
	IPv6Addresses []string

	// IPv6 Prefixes allocated for the network interface
	IPv6Prefixes []string
}

// ReconcileCooldownCache keep track of recently freed CIDRs to avoid reading stale EC2 metadata
type ReconcileCooldownCache struct {
	sync.RWMutex
	cache map[string]time.Time
}

// Add sets a timestamp for the CIDR added that says how long they are not to be put back in the data store.
func (r *ReconcileCooldownCache) Add(cidr string) {
	r.Lock()
	defer r.Unlock()
	expiry := time.Now().Add(ipReconcileCooldown)
	r.cache[cidr] = expiry
}

// Remove removes a CIDR from the cooldown cache.
func (r *ReconcileCooldownCache) Remove(cidr string) {
	r.Lock()
	defer r.Unlock()
	log.Debugf("Removing %s from cooldown cache.", cidr)
	delete(r.cache, cidr)
}

// RecentlyFreed checks if this CIDR was recently freed.
func (r *ReconcileCooldownCache) RecentlyFreed(cidr string) (found, recentlyFreed bool) {
	r.Lock()
	defer r.Unlock()
	now := time.Now()
	if expiry, ok := r.cache[cidr]; ok {
		log.Debugf("Checking if CIDR %s has been recently freed. Cooldown expires at: %s. (Cooldown: %v)", cidr, expiry, now.Sub(expiry) < 0)
		return true, now.Sub(expiry) < 0
	}
	return false, false
}

func prometheusRegister() {
	if !prometheusRegistered {
		prometheusmetrics.PrometheusRegister()
		prometheusRegistered = true
	}
}

// New retrieves IP address usage information from Instance MetaData service and Kubelet
// then initializes IP address pool data store
func New(k8sClient client.Client, enableIPv4, enableIPv6 bool, nodeName string, v4VPCCIDRs, v6VPCCIDRs []string,
	primaryIP, primaryMAC, primaryENIID string, ipDataStore *datastore.DataStore, networkInterfaces []ENIMetadata) (*IPAMContext, error) {
	var err error

	prometheusRegister()
	c := &IPAMContext{}
	c.k8sClient = k8sClient
	c.networkClient = networkutils.New()
	c.enableIPv4 = enableIPv4
	c.enableIPv6 = enableIPv6
	c.v4VPCCIDRs = v4VPCCIDRs
	c.v6VPCCIDRs = v6VPCCIDRs
	c.dataStore = ipDataStore

	c.primaryIP = make(map[string]string)
	//c.enablePodIPAnnotation = enablePodIPAnnotation()
	//c.networkPolicyMode, err = getNetworkPolicyMode()
	if err != nil {
		return nil, err
	}

	c.myNodeName = nodeName
	c.reconcileCooldownCache.cache = make(map[string]time.Time)

	if err := c.nodeInit(v4VPCCIDRs, primaryMAC, primaryIP, primaryENIID, networkInterfaces); err != nil {
		return nil, err
	}
	return c, nil
}

// TODO - Call NodeInit from CNINode reconciler
func (c *IPAMContext) nodeInit(vpcV4CIDRs []string, primaryENIMac, primaryIP, primaryENIID string, enis []ENIMetadata) error {
	prometheusmetrics.IpamdActionsInprogress.WithLabelValues("nodeInit").Add(float64(1))
	defer prometheusmetrics.IpamdActionsInprogress.WithLabelValues("nodeInit").Sub(float64(1))
	var err error
	//ctx := context.TODO()

	log.Debugf("Start node init")
	primaryV4IP := net.ParseIP(primaryIP)
	err = c.networkClient.SetupHostNetwork(vpcV4CIDRs, primaryENIMac, &primaryV4IP, false, true, false)

	if err != nil {
		return errors.Wrap(err, "ipamd init: failed to set up host network")
	}
	err = c.networkClient.CleanUpStaleAWSChains(c.enableIPv4, c.enableIPv6)
	if err != nil {
		// We should not error if clean up fails since these chains don't affect the rules
		log.Debugf("Failed to clean up stale AWS chains: %v", err)
	}

	for _, eni := range enis {
		log.Debugf("Discovered ENI %s, trying to set it up", eni.ENIID)
		///isEFAENI := metadataResult.EFAENIs[eni.ENIID]

		// Retry ENI sync
		retry := 0
		for {
			retry++
			if err = c.SetupENI(eni, false, false, primaryENIID); err == nil {
				log.Infof("ENI %s set up.", eni.ENIID)
				break
			}

			if retry > maxRetryCheckENI {
				log.Warnf("Reached max retry: Unable to discover attached IPs for ENI from metadata service (attempted %d/%d): %v", retry, maxRetryCheckENI, err)
				ipamdErrInc("waitENIAttachedMaxRetryExceeded")
				break
			}

			log.Warnf("Error trying to set up ENI %s: %v", eni.ENIID, err)
			if strings.Contains(err.Error(), "setupENINetwork: failed to find the link which uses MAC address") {
				// If we can't find the matching link for this MAC address, there is no point in retrying for this ENI.
				log.Debug("Unable to match link for this ENI, going to the next one.")
				break
			}
			log.Debugf("Unable to discover IPs for this ENI yet (attempt %d/%d)", retry, maxRetryCheckENI)
			time.Sleep(eniAttachTime)
		}
	}

	if err := c.dataStore.ReadBackingStore(c.enableIPv6); err != nil {
		return err
	}

	/*
		if err = c.configureIPRulesForPods(); err != nil {
			return err
		}
		// Spawning updateCIDRsRulesOnChange go-routine
		go wait.Forever(func() {
			vpcV4CIDRs = c.updateCIDRsRulesOnChange(vpcV4CIDRs)
		}, 30*time.Second)

	*/

	log.Debug("node init completed successfully")
	return nil
}

func (c *IPAMContext) configureIPRulesForPods() error {
	rules, err := c.networkClient.GetRuleList()
	if err != nil {
		log.Errorf("During ipamd init: failed to retrieve IP rule list %v", err)
		return nil
	}

	for _, info := range c.dataStore.AllocatedIPs() {
		// TODO(gus): This should really be done via CNI CHECK calls, rather than in ipam (requires upstream k8s changes).

		// Update ip rules in case there is a change in VPC CIDRs, AWS_VPC_K8S_CNI_EXTERNALSNAT setting
		srcIPNet := net.IPNet{IP: net.ParseIP(info.IP), Mask: net.IPv4Mask(255, 255, 255, 255)}

		err = c.networkClient.UpdateRuleListBySrc(rules, srcIPNet)
		if err != nil {
			log.Warnf("UpdateRuleListBySrc in nodeInit() failed for IP %s: %v", info.IP, err)
		}
	}

	// Program IP rules for external service CIDRs and cleanup stale rules.
	// Note that we can reuse rule list despite it being modified by UpdateRuleListBySrc, as the
	// modifications touched rules that this function ignores.
	extServiceCIDRs := c.networkClient.GetExternalServiceCIDRs()
	err = c.networkClient.UpdateExternalServiceIpRules(rules, extServiceCIDRs)
	if err != nil {
		log.Warnf("UpdateExternalServiceIpRules in nodeInit() failed")
	}

	return nil
}

//TODO - Call this function if we detect a change in VPCCIDRs (or) maybe let's depend on IMDS for it..
/*
func (c *IPAMContext) updateCIDRsRulesOnChange(oldVPCCIDRs []string) []string {
	newVPCCIDRs, err := c.awsClient.GetVPCIPv4CIDRs()
	if err != nil {
		log.Warnf("skipping periodic update to VPC CIDRs due to error: %v", err)
		return oldVPCCIDRs
	}

	old := sets.NewString(oldVPCCIDRs...)
	new := sets.NewString(newVPCCIDRs...)
	if !old.Equal(new) {
		primaryIP := c.awsClient.GetLocalIPv4()
		err = c.networkClient.UpdateHostIptablesRules(newVPCCIDRs, c.awsClient.GetPrimaryENImac(), &primaryIP, c.enableIPv4,
			c.enableIPv6)
		if err != nil {
			log.Warnf("unable to update host iptables rules for VPC CIDRs due to error: %v", err)
		}
	}
	return newVPCCIDRs
}
*/

//TODO - REQUIRED
// setupENI does following:
// 1) add ENI to datastore
// 2) set up linux ENI related networking stack.
// 3) add all ENI's secondary IP addresses to datastore

func (c *IPAMContext) SetupENI(eniMetadata ENIMetadata, isTrunkENI, isEFAENI bool, primaryENI string) error {
	// Add the ENI to the datastore
	err := c.dataStore.AddENI(eniMetadata.ENIID, int(eniMetadata.DeviceNumber), eniMetadata.ENIID == primaryENI, isTrunkENI, isEFAENI)
	if err != nil && err.Error() != datastore.DuplicatedENIError {
		return errors.Wrapf(err, "failed to add ENI %s to data store", eniMetadata.ENIID)
	}
	// Store the addressable IP for the ENI
	//TODO - Fill the Primary ENI
	/*
		if c.enableIPv6 {
			c.primaryIP[eni] = eniMetadata.PrimaryIPv6Address()
		} else {
			c.primaryIP[eni] = eniMetadata.PrimaryIPv4Address()
		}
	*/

	c.primaryIP[eniMetadata.ENIID] = eniMetadata.PrimaryIP
	// In v6 PD mode, VPC CNI will only manage the primary ENI and trunk ENI. Once we start supporting secondary
	// IP and custom networking modes for IPv6, this restriction can be relaxed.
	log.Infof("ENI ID: %s; Primary ENI ID: %s", eniMetadata.ENIID, primaryENI)

	//TODO - Handle IPv6 Mode
	// For other ENIs, set up the network
	subnetCidr := eniMetadata.SubnetIPv4CIDR
	if c.enableIPv6 {
		subnetCidr = eniMetadata.SubnetIPv6CIDR
	}
	if eniMetadata.ENIID != primaryENI {
		log.Infof("Secondary ENI. Let's configure it....")
		err = c.networkClient.SetupENINetwork(c.primaryIP[eniMetadata.ENIID], eniMetadata.MAC, int(eniMetadata.DeviceNumber), subnetCidr)
		if err != nil {
			// Failed to set up the ENI
			errRemove := c.dataStore.RemoveENIFromDataStore(eniMetadata.ENIID, true)
			if errRemove != nil {
				log.Warnf("failed to remove ENI %s: %v", eniMetadata.ENIID, errRemove)
			}
			delete(c.primaryIP, eniMetadata.ENIID)
			return errors.Wrapf(err, "failed to set up ENI %s network", eniMetadata.ENIID)
		}
	}

	log.Infof("Found ENIs having %d secondary IPs and %d Prefixes", len(eniMetadata.IPv4Addresses), len(eniMetadata.IPv4Prefixes))
	// Either case add the IPs and prefixes to datastore.
	c.addENIsecondaryIPsToDataStore(eniMetadata.IPv4Addresses, eniMetadata.ENIID, eniMetadata.PrimaryIP)
	c.addENIv4prefixesToDataStore(eniMetadata.IPv4Prefixes, eniMetadata.ENIID)

	return nil
}

func (c *IPAMContext) addENIsecondaryIPsToDataStore(PrivateIpAddrs []string, eni string, primaryIP string) {
	// Add all the secondary IPs
	for _, PrivateIpAddress := range PrivateIpAddrs {
		log.Infof("IP Address: %s ; Primary IP: %s", PrivateIpAddress, primaryIP)
		if strings.Compare(PrivateIpAddress+hostMask, primaryIP) == 0 {
			log.Infof("Skip adding Primary IP to datastore....")
			continue
		}
		cidr := net.IPNet{IP: net.ParseIP(PrivateIpAddress), Mask: net.IPv4Mask(255, 255, 255, 255)}
		err := c.dataStore.AddIPv4CidrToStore(eni, cidr, false)
		if err != nil && err.Error() != datastore.IPAlreadyInStoreError {
			log.Warnf("Failed to increase IP pool, failed to add IP %s to data store", PrivateIpAddress)
			// continue to add next address
			ipamdErrInc("addENIsecondaryIPsToDataStoreFailed")
		}
	}
	//c.logPoolStats(c.dataStore.GetIPStats(ipV4AddrFamily))

}

func (c *IPAMContext) addENIv4prefixesToDataStore(PrefixAddrs []string, eni string) {
	// Walk thru all prefixes
	for _, PrefixAddr := range PrefixAddrs {
		log.Infof("Adding Prefix: %s", PrefixAddr)
		strIpv4Prefix := PrefixAddr
		_, ipnet, err := net.ParseCIDR(strIpv4Prefix)
		if err != nil {
			//Parsing failed, get next prefix
			log.Debugf("Parsing failed, moving on to next prefix")
			continue
		}
		cidr := *ipnet
		err = c.dataStore.AddIPv4CidrToStore(eni, cidr, true)
		if err != nil && err.Error() != datastore.IPAlreadyInStoreError {
			log.Warnf("Failed to increase Prefix pool, failed to add Prefix %s to data store", PrefixAddr)
			// continue to add next address
			ipamdErrInc("addENIv4prefixesToDataStoreFailed")
		}
	}
	c.logPoolStats(c.dataStore.GetIPStats(ipV4AddrFamily))
}

func (c *IPAMContext) addENIv6prefixesToDataStore(PrefixAddrs []string, eni string) {
	// Walk through all prefixes
	for _, PrefixAddr := range PrefixAddrs {
		strIpv6Prefix := PrefixAddr
		_, ipnet, err := net.ParseCIDR(strIpv6Prefix)
		if err != nil {
			// Parsing failed, get next prefix
			log.Debugf("Parsing failed, moving on to next prefix")
			continue
		}
		cidr := *ipnet
		err = c.dataStore.AddIPv6CidrToStore(eni, cidr, true)
		if err != nil && err.Error() != datastore.IPAlreadyInStoreError {
			log.Warnf("Failed to increase Prefix pool, failed to add Prefix %s to data store", PrefixAddr)
			// continue to add next address
			ipamdErrInc("addENIv6prefixesToDataStoreFailed")
		}
	}
	c.logPoolStats(c.dataStore.GetIPStats(ipV6AddrFamily))
}

// logPoolStats logs usage information for allocated addresses/prefixes.
func (c *IPAMContext) logPoolStats(dataStoreStats *datastore.DataStoreStats) {
	prefix := "IP pool stats"
	log.Debugf("%s: %s", prefix, dataStoreStats)
}

func ipamdErrInc(fn string) {
	prometheusmetrics.IpamdErr.With(prometheus.Labels{"fn": fn}).Inc()
}

func podENIErrInc(fn string) {
	prometheusmetrics.PodENIErr.With(prometheus.Labels{"fn": fn}).Inc()
}

// verifyAndAddIPsToDatastore updates the datastore with the known secondary IPs. IPs who are out of cooldown gets added
// back to the datastore after being verified against EC2.
func (c *IPAMContext) VerifyAndAddIPsToDatastore(eni string, attachedENIIPs []string) map[string]bool {
	seenIPs := make(map[string]bool)
	for _, privateIPv4 := range attachedENIIPs {
		strPrivateIPv4 := privateIPv4 //aws.StringValue(privateIPv4.PrivateIpAddress)
		if strPrivateIPv4 == c.primaryIP[eni] {
			log.Infof("Reconcile and skip primary IP %s on ENI %s", strPrivateIPv4, eni)
			continue
		}

		// Check if this IP was recently freed
		ipv4Addr := net.IPNet{IP: net.ParseIP(strPrivateIPv4), Mask: net.IPv4Mask(255, 255, 255, 255)}
		found, recentlyFreed := c.reconcileCooldownCache.RecentlyFreed(strPrivateIPv4)
		if found {
			if recentlyFreed {
				log.Debugf("Reconcile skipping IP %s on ENI %s because it was recently unassigned from the ENI.", strPrivateIPv4, eni)
				continue
			} else {
				// The IP can be removed from the cooldown cache
				c.reconcileCooldownCache.Remove(strPrivateIPv4)
			}
		}
		log.Infof("Trying to add %s", strPrivateIPv4)
		// Try to add the IP
		err := c.dataStore.AddIPv4CidrToStore(eni, ipv4Addr, false)
		if err != nil && err.Error() != datastore.IPAlreadyInStoreError {
			log.Errorf("Failed to reconcile IP %s on ENI %s", strPrivateIPv4, eni)
			ipamdErrInc("ipReconcileAdd")
			// Continue to check the other IPs instead of bailout due to one wrong IP
			continue

		}
		// Mark action
		seenIPs[strPrivateIPv4] = true
		prometheusmetrics.ReconcileCnt.With(prometheus.Labels{"fn": "eniDataStorePoolReconcileAdd"}).Inc()
	}
	return seenIPs
}

// verifyAndAddPrefixesToDatastore updates the datastore with the known Prefixes. Prefixes who are out of cooldown gets added
// back to the datastore after being verified against EC2.
func (c *IPAMContext) VerifyAndAddPrefixesToDatastore(eni string, attachedENIPrefixes []string) map[string]bool {
	seenIPs := make(map[string]bool)
	for _, privateIPv4Cidr := range attachedENIPrefixes {
		strPrivateIPv4Cidr := privateIPv4Cidr
		log.Debugf("Check in coolddown Found prefix %s", strPrivateIPv4Cidr)

		// Check if this Prefix was recently freed
		_, ipv4CidrPtr, err := net.ParseCIDR(strPrivateIPv4Cidr)
		if err != nil {
			log.Debugf("Failed to parse so continuing with next prefix")
			continue
		}
		found, recentlyFreed := c.reconcileCooldownCache.RecentlyFreed(strPrivateIPv4Cidr)
		if found {
			if recentlyFreed {
				log.Debugf("Reconcile skipping IP %s on ENI %s because it was recently unassigned from the ENI.", strPrivateIPv4Cidr, eni)
				continue
			} else {
				// The IP can be removed from the cooldown cache
				// TODO: Here we could check if the Prefix is still used by a pod stuck in Terminating state. (Issue #1091)
				c.reconcileCooldownCache.Remove(strPrivateIPv4Cidr)
			}
		}

		err = c.dataStore.AddIPv4CidrToStore(eni, *ipv4CidrPtr, true)
		if err != nil && err.Error() != datastore.IPAlreadyInStoreError {
			log.Errorf("Failed to reconcile Prefix %s on ENI %s", strPrivateIPv4Cidr, eni)
			ipamdErrInc("prefixReconcileAdd")
			// Continue to check the other Prefixs instead of bailout due to one wrong IP
			continue

		}
		// Mark action
		seenIPs[strPrivateIPv4Cidr] = true
		prometheusmetrics.ReconcileCnt.With(prometheus.Labels{"fn": "eniDataStorePoolReconcileAdd"}).Inc()
	}
	return seenIPs
}

func parseBoolEnvVar(envVariableName string, defaultVal bool) bool {
	if strValue := os.Getenv(envVariableName); strValue != "" {
		parsedValue, err := strconv.ParseBool(strValue)
		if err == nil {
			return parsedValue
		}
		log.Warnf("Failed to parse %s; using default: %v, err: %v", envVariableName, defaultVal, err)
	}
	return defaultVal
}

func disableENIProvisioning() bool {
	return utils.GetBoolAsStringEnvVar(envDisableENIProvisioning, false)
}

// setTerminating atomically sets the terminating flag.
func (c *IPAMContext) setTerminating() {
	atomic.StoreInt32(&c.terminating, 1)
}

func (c *IPAMContext) isTerminating() bool {
	return atomic.LoadInt32(&c.terminating) > 0
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func min(x, y int) int {
	if y < x {
		return y
	}
	return x
}
