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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type CIDR string

type NetworkInterface struct {
	// +kubebuilder:validation:XValidation:rule="self != ''",message="must not be empty"
	ID string `json:"id"`
	// +kubebuilder:validation:XValidation:rule="self != ''",message="must not be empty"
	SubnetID string `json:"subnetId"`
	// +kubebuilder:validation:XValidation:rule="self != ''",message="must not be empty"
	AttachmentId string `json:"attachmentId"`
	// +kubebuilder:validation:XValidation:rule="self >= 0",message="must not be negative"
	NetworkCardIndex int32 `json:"networkCardIndex"`
	// +kubebuilder:validation:XValidation:rule="self >= 0",message="must not be negative"
	DeviceIndex int32 `json:"deviceIndex"`
	// +kubebuilder:validation:XValidation:rule="self != ''",message="must not be empty"
	PrimaryCIDR CIDR `json:"primaryCIDR"`
	// +kubebuilder:validation:XValidation:message="must not be empty",rule="self.size() != 0"
	CIDRs []CIDR `json:"cidrs"`
}

// CNINodeSpec defines the desired state of CNINode
type CNINodeSpec struct {
	// +kubebuilder:validation:XValidation:rule="self != ''",message="must not be empty"
	InstanceID string `json:"instanceID"`
	// +kubebuilder:validation:XValidation:rule="self != ''",message="must not be empty"
	InstanceType types.InstanceType `json:"instanceType"`
	// +kubebuilder:validation:XValidation:rule="self != ''",message="must not be empty"
	Hypervisor types.InstanceTypeHypervisor `json:"hypervisor"`
	// +kubebuilder:validation:XValidation:rule="self != ''",message="must not be empty"
	PrimaryIPv4Address string `json:"primaryIPv4"`
	// +kubebuilder:validation:XValidation:rule="self != ''",message="must not be empty"
	PrimaryNetworkInterfaceID string `json:"primaryNetworkInterfaceID"`
	// +kubebuilder:validation:XValidation:rule="self > 0",message="must be positive"
	MaxiumumNetworkCards int32 `json:"maximumNetworkCards"`
	// +kubebuilder:validation:XValidation:rule="self > 0",message="must be positive"
	MaximumNetworkInterfaces int32 `json:"maximumNetworkInterfaces"`
	// +kubebuilder:validation:XValidation:rule="self > 0",message="must be positive"
	Ipv4AddressesPerInterface int32 `json:"ipv4AddressesPerInterface"`
}

// CNINodeStatus defines the observed state of CNINode
type CNINodeStatus struct {
	NetworkInterfaces []NetworkInterface `json:"networkInterfaces,omitempty"`
	SubnetIDs         []string           `json:"subnetIDs,omitempty"`
	VPCCIDRs          []string           `json:"vpcCIDRs,omitempty"`
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
