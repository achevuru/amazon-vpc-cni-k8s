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
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"

	"github.com/aws/amazon-vpc-cni-k8s/pkg/ipamd/datastore"
	"github.com/aws/amazon-vpc-cni-k8s/rpc"
	"github.com/aws/amazon-vpc-cni-k8s/utils/prometheusmetrics"
)

const (
	ipamdgRPCaddress      = "127.0.0.1:50051"
	grpcHealthServiceName = "grpc.health.v1.aws-node"

	vpccniPodIPKey = "vpc.amazonaws.com/pod-ips"
)

// server controls RPC service responses.
type server struct {
	version     string
	ipamContext *IPAMContext
}

// PodENIData is used to parse the list of ENIs in the branch ENI pod annotation
type PodENIData struct {
	ENIID        string `json:"eniId"`
	IfAddress    string `json:"ifAddress"`
	PrivateIP    string `json:"privateIp"`
	IPV6Addr     string `json:"ipv6Addr"`
	VlanID       int    `json:"vlanID"`
	SubnetCIDR   string `json:"subnetCidr"`
	SubnetV6CIDR string `json:"subnetV6Cidr"`
}

// AddNetwork processes CNI add network request and return an IP address for container
func (s *server) AddNetwork(ctx context.Context, in *rpc.AddNetworkRequest) (*rpc.AddNetworkReply, error) {
	log.Infof("Received AddNetwork for NS %s, Sandbox %s, ifname %s",
		in.Netns, in.ContainerID, in.IfName)
	log.Debugf("AddNetworkRequest: %s", in)
	prometheusmetrics.AddIPCnt.Inc()

	// Do this early, but after logging trace
	if err := s.validateVersion(in.ClientVersion); err != nil {
		log.Warnf("Rejecting AddNetwork request: %v", err)
		return nil, err
	}

	failureResponse := rpc.AddNetworkReply{Success: false}
	var deviceNumber, vlanID, trunkENILinkIndex int
	var ipv4Addr, ipv6Addr, branchENIMAC, podENISubnetGW string
	var err error

	if s.ipamContext.enableIPv4 && ipv4Addr == "" ||
		s.ipamContext.enableIPv6 && ipv6Addr == "" {
		if in.ContainerID == "" || in.IfName == "" || in.NetworkName == "" {
			log.Errorf("Unable to generate IPAMKey from %+v", in)
			return &failureResponse, nil
		}
		ipamKey := datastore.IPAMKey{
			ContainerID: in.ContainerID,
			IfName:      in.IfName,
			NetworkName: in.NetworkName,
		}
		ipamMetadata := datastore.IPAMMetadata{
			K8SPodNamespace: in.K8S_POD_NAMESPACE,
			K8SPodName:      in.K8S_POD_NAME,
		}
		ipv4Addr, ipv6Addr, deviceNumber, err = s.ipamContext.dataStore.AssignPodIPAddress(ipamKey, ipamMetadata, s.ipamContext.enableIPv4, s.ipamContext.enableIPv6)
	}

	var pbVPCV4cidrs, pbVPCV6cidrs []string
	var useExternalSNAT bool
	if s.ipamContext.enableIPv4 && ipv4Addr != "" {
		pbVPCV4cidrs = s.ipamContext.v4VPCCIDRs
		for _, cidr := range pbVPCV4cidrs {
			log.Debugf("VPC CIDR %s", cidr)
		}
		useExternalSNAT = s.ipamContext.isExternalSNATEnabled()
		if !useExternalSNAT {
			for _, cidr := range s.ipamContext.networkClient.GetExcludeSNATCIDRs() {
				log.Debugf("CIDR SNAT Exclusion %s", cidr)
				pbVPCV4cidrs = append(pbVPCV4cidrs, cidr)
			}
		}
	} else if s.ipamContext.enableIPv6 && ipv6Addr != "" {
		pbVPCV6cidrs = s.ipamContext.v6VPCCIDRs
		for _, cidr := range pbVPCV6cidrs {
			log.Debugf("VPC V6 CIDR %s", cidr)
		}
	}

	/*
		if s.ipamContext.enablePodIPAnnotation {
			// On ADD, we pass empty string as there is no IP being released
			if ipv4Addr != "" {
				err = s.ipamContext.AnnotatePod(in.K8S_POD_NAME, in.K8S_POD_NAMESPACE, vpccniPodIPKey, ipv4Addr, "")
				if err != nil {
					log.Errorf("Failed to add the pod annotation: %v", err)
				}
			} else if ipv6Addr != "" {
				err = s.ipamContext.AnnotatePod(in.K8S_POD_NAME, in.K8S_POD_NAMESPACE, vpccniPodIPKey, ipv6Addr, "")
				if err != nil {
					log.Errorf("Failed to add the pod annotation: %v", err)
				}
			}
		}
	*/
	resp := rpc.AddNetworkReply{
		Success:           err == nil,
		IPv4Addr:          ipv4Addr,
		IPv6Addr:          ipv6Addr,
		DeviceNumber:      int32(deviceNumber),
		UseExternalSNAT:   useExternalSNAT,
		VPCv4CIDRs:        pbVPCV4cidrs,
		VPCv6CIDRs:        pbVPCV6cidrs,
		PodVlanId:         int32(vlanID),
		PodENIMAC:         branchENIMAC,
		PodENISubnetGW:    podENISubnetGW,
		ParentIfIndex:     int32(trunkENILinkIndex),
		NetworkPolicyMode: s.ipamContext.networkPolicyMode,
	}

	log.Infof("Send AddNetworkReply: IPv4Addr: %s, IPv6Addr: %s, DeviceNumber: %d, err: %v", ipv4Addr, ipv6Addr, deviceNumber, err)
	return &resp, nil
}

func (s *server) validateVersion(clientVersion string) error {
	if s.version != clientVersion {
		return status.Errorf(codes.FailedPrecondition, "wrong client version %q (!= %q)", clientVersion, s.version)
	}
	return nil
}

func (s *server) DelNetwork(ctx context.Context, in *rpc.DelNetworkRequest) (*rpc.DelNetworkReply, error) {
	log.Infof("Received DelNetwork for Sandbox %s", in.ContainerID)
	log.Debugf("DelNetworkRequest: %s", in)
	prometheusmetrics.DelIPCnt.With(prometheus.Labels{"reason": in.Reason}).Inc()
	var ipv4Addr, ipv6Addr string

	// Do this early, but after logging trace
	if err := s.validateVersion(in.ClientVersion); err != nil {
		log.Warnf("Rejecting DelNetwork request: %v", err)
		return nil, err
	}

	ipamKey := datastore.IPAMKey{
		ContainerID: in.ContainerID,
		IfName:      in.IfName,
		NetworkName: in.NetworkName,
	}
	_, ip, deviceNumber, err := s.ipamContext.dataStore.UnassignPodIPAddress(ipamKey)
	if s.ipamContext.enableIPv4 {
		ipv4Addr = ip
	} else if s.ipamContext.enableIPv6 {
		ipv6Addr = ip
	}

	/*
		if err == datastore.ErrUnknownPod && s.ipamContext.enablePodENI {
			pod, err := s.ipamContext.GetPod(in.K8S_POD_NAME, in.K8S_POD_NAMESPACE)
			if err != nil {
				if k8serror.IsNotFound(err) {
					log.Warn("Send DelNetworkReply: pod not found")
					return &rpc.DelNetworkReply{Success: true}, nil
				}
				log.Warnf("Send DelNetworkReply: Failed to get pod spec: %v", err)
				return &rpc.DelNetworkReply{Success: false}, err
			}
			val, branch := pod.Annotations["vpc.amazonaws.com/pod-eni"]
			if branch {
				// Parse JSON data
				var podENIData []PodENIData
				err := json.Unmarshal([]byte(val), &podENIData)
				if err != nil || len(podENIData) < 1 {
					log.Errorf("Failed to unmarshal PodENIData JSON: %v", err)
				}
				return &rpc.DelNetworkReply{
					Success:   true,
					PodVlanId: int32(podENIData[0].VlanID),
					IPv4Addr:  podENIData[0].PrivateIP}, err
			}
		}
	*/

	/*
		if s.ipamContext.enablePodIPAnnotation {
			// On DEL, we pass IP being released
			err = s.ipamContext.AnnotatePod(in.K8S_POD_NAME, in.K8S_POD_NAMESPACE, vpccniPodIPKey, "", ip)
			if err != nil {
				log.Errorf("Failed to delete the pod annotation: %v", err)
			}
		}
	*/

	log.Infof("Send DelNetworkReply: IPv4Addr: %s, IPv6Addr: %s, DeviceNumber: %d, err: %v", ipv4Addr, ipv6Addr, deviceNumber, err)

	return &rpc.DelNetworkReply{Success: err == nil, IPv4Addr: ipv4Addr, IPv6Addr: ipv6Addr, DeviceNumber: int32(deviceNumber)}, err
}

// RunRPCHandler handles request from gRPC
func (c *IPAMContext) RunRPCHandler(version string) error {
	log.Infof("Serving RPC Handler version %s on %s", version, ipamdgRPCaddress)
	listener, err := net.Listen("tcp", ipamdgRPCaddress)
	if err != nil {
		log.Errorf("Failed to listen gRPC port: %v", err)
		return errors.Wrap(err, "ipamd: failed to listen to gRPC port")
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterCNIBackendServer(grpcServer, &server{version: version, ipamContext: c})
	healthServer := health.NewServer()
	// If ipamd can talk to the API server and to the EC2 API, the pod is healthy.
	// No need to ever change this to HealthCheckResponse_NOT_SERVING since it's a local service only
	healthServer.SetServingStatus(grpcHealthServiceName, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Errorf("Failed to start server on gRPC port: %v", err)
		return errors.Wrap(err, "ipamd: failed to start server on gPRC port")
	}
	return nil
}
