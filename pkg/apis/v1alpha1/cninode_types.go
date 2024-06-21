/*
Copyright 2024.

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

package v1alpha1

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/operatorpkg/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type CIDR string

type UnusedCIDRs struct {
	CIDRs           []CIDR      `json:"cidrs"`
	UnusedTimestamp metav1.Time `json:"unusedTimestamp"`
}

type IPAllocationStrategy string

const (
	IPAllocationStrategyPrefixFallback IPAllocationStrategy = "PrefixFallback"
	IPAllocationStrategySecondaryIP    IPAllocationStrategy = "SecondaryIP"
)

type NetworkPolicyEnforcement string

const (
	StrictMode NetworkPolicyEnforcement = "DefaultDeny"
	Permissive NetworkPolicyEnforcement = "DefaultAllow"
)

type SNATPolicy string

const (
	Hashrandom = "hashrandom"
	PRNG       = "prng"
	None       = "none"
	External   = "external"
)

type NetworkInterface struct {
	// +required
	// +kubebuilder:validation:MinLength:=1
	ID string `json:"id"`
	// +required
	// +kubebuilder:validation:MinLength:=1
	SubnetID string `json:"subnetId"`
	// +required
	// +kubebuilder:validation:MinLength:=1
	AttachmentId string `json:"attachmentId"`
	// +required
	// +kubebuilder:validation:Minimum:=0
	NetworkCardIndex int32 `json:"networkCardIndex"`
	// +required
	// +kubebuilder:validation:Minimum:=0
	DeviceIndex int32 `json:"deviceIndex"`
	// +required
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:Type:=string
	PrimaryCIDR CIDR `json:"primaryCIDR"`
	// +required
	// +kubebuilder:validation:MinLength:=1
	MacAddress string `json:"macAddress"`
	// +required
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:Type:=string
	SubnetV4CIDR CIDR `json:"subnetV4CIDR"`
	SubnetV6CIDR CIDR `json:"subnetV6CIDR"`
	// +required
	// +kubebuilder:validation:MinItems:=1
	V4CIDRs []CIDR `json:"v4CIDRs"`
	// +optional
	V6CIDRs []CIDR `json:"v6CIDRs"`
	// +optional
	UnusedCIDRs []UnusedCIDRs `json:"unusedCIDRs"`
}

// CNINodeSpec defines the desired state of CNINode
type CNINodeSpec struct {
	// +required
	// +kubebuilder:validation:MinLength:=1
	InstanceID string `json:"instanceID"`
	// +required
	// +kubebuilder:validation:MinLength:=1
	VPCID string `json:"vpcID"`
	// +required
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:Type:=string
	InstanceType types.InstanceType `json:"instanceType"`
	// +required
	// +kubebuilder:validation:Enum:={IPv4,IPv6}
	IPFamily corev1.IPFamily `json:"ipFamily"`
	// +required
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:Type:=string
	IPAllocationStrategy IPAllocationStrategy `json:"ipAllocationStrategy"`
	// +required
	// +kubebuilder:validation:Enum=hashrandom;prng;none;external
	SNATPolicy SNATPolicy `json:"snatPolicy"`
	// +required
	// +kubebuilder:validation:Enum=DefaultDeny;DefaultAllow
	NetworkPolicyEnforcement NetworkPolicyEnforcement `json:"networkPolicyEnforcement"`
	// +required
	// +kubebuilder:validation:Minimum:=1
	ConntrackCleanupInterval int `json:"conntrackCleanupInterval"`
	// +required
	// +kubebuilder:validation:MinLength:=1
	PrimaryIPv4Address string `json:"primaryIPv4"`
	// +required
	// +kubebuilder:validation:MinLength:=1
	PrimaryNetworkInterfaceID string `json:"primaryNetworkInterfaceID"`
	// +required
	// +kubebuilder:validation:Minimum:=1
	MaximumNetworkCards int32 `json:"maximumNetworkCards"`
	// +required
	// +kubebuilder:validation:Minimum:=1
	MaximumNetworkInterfaces int32 `json:"maximumNetworkInterfaces"`
	// +required
	// +kubebuilder:validation:Minimum:=1
	IPv4AddressesPerInterface int32 `json:"ipv4AddressesPerInterface"`
}

// CNINodeStatus defines the observed state of CNINode
type CNINodeStatus struct {
	NetworkInterfaces []NetworkInterface `json:"networkInterfaces,omitempty"`
	SubnetIDs         []string           `json:"subnetIDs,omitempty"`
	VPCCIDRs          []string           `json:"vpcCIDRs,omitempty"`
	VPCV6CIDRs        []string           `json:"vpcV6CIDRs,omitempty"`
	Conditions        []status.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CNINode is the Schema for the cninodes API
type CNINode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CNINodeSpec   `json:"spec,omitempty"`
	Status CNINodeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CNINodeList contains a list of CNINode
type CNINodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CNINode `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CNINode{}, &CNINodeList{})
}
