// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.2
// source: rpc/rpc.proto

package rpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AddNetworkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientVersion              string `protobuf:"bytes,8,opt,name=ClientVersion,proto3" json:"ClientVersion,omitempty"`
	K8S_POD_NAME               string `protobuf:"bytes,1,opt,name=K8S_POD_NAME,json=K8SPODNAME,proto3" json:"K8S_POD_NAME,omitempty"`
	K8S_POD_NAMESPACE          string `protobuf:"bytes,2,opt,name=K8S_POD_NAMESPACE,json=K8SPODNAMESPACE,proto3" json:"K8S_POD_NAMESPACE,omitempty"`
	K8S_POD_INFRA_CONTAINER_ID string `protobuf:"bytes,3,opt,name=K8S_POD_INFRA_CONTAINER_ID,json=K8SPODINFRACONTAINERID,proto3" json:"K8S_POD_INFRA_CONTAINER_ID,omitempty"`
	ContainerID                string `protobuf:"bytes,7,opt,name=ContainerID,proto3" json:"ContainerID,omitempty"`
	IfName                     string `protobuf:"bytes,5,opt,name=IfName,proto3" json:"IfName,omitempty"`
	NetworkName                string `protobuf:"bytes,6,opt,name=NetworkName,proto3" json:"NetworkName,omitempty"`
	Netns                      string `protobuf:"bytes,4,opt,name=Netns,proto3" json:"Netns,omitempty"` // next field: 9
}

func (x *AddNetworkRequest) Reset() {
	*x = AddNetworkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddNetworkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddNetworkRequest) ProtoMessage() {}

func (x *AddNetworkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddNetworkRequest.ProtoReflect.Descriptor instead.
func (*AddNetworkRequest) Descriptor() ([]byte, []int) {
	return file_rpc_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *AddNetworkRequest) GetClientVersion() string {
	if x != nil {
		return x.ClientVersion
	}
	return ""
}

func (x *AddNetworkRequest) GetK8S_POD_NAME() string {
	if x != nil {
		return x.K8S_POD_NAME
	}
	return ""
}

func (x *AddNetworkRequest) GetK8S_POD_NAMESPACE() string {
	if x != nil {
		return x.K8S_POD_NAMESPACE
	}
	return ""
}

func (x *AddNetworkRequest) GetK8S_POD_INFRA_CONTAINER_ID() string {
	if x != nil {
		return x.K8S_POD_INFRA_CONTAINER_ID
	}
	return ""
}

func (x *AddNetworkRequest) GetContainerID() string {
	if x != nil {
		return x.ContainerID
	}
	return ""
}

func (x *AddNetworkRequest) GetIfName() string {
	if x != nil {
		return x.IfName
	}
	return ""
}

func (x *AddNetworkRequest) GetNetworkName() string {
	if x != nil {
		return x.NetworkName
	}
	return ""
}

func (x *AddNetworkRequest) GetNetns() string {
	if x != nil {
		return x.Netns
	}
	return ""
}

type AddNetworkReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success         bool     `protobuf:"varint,1,opt,name=Success,proto3" json:"Success,omitempty"`
	IPv4Addr        string   `protobuf:"bytes,2,opt,name=IPv4Addr,proto3" json:"IPv4Addr,omitempty"`
	IPv6Addr        string   `protobuf:"bytes,3,opt,name=IPv6Addr,proto3" json:"IPv6Addr,omitempty"`
	DeviceNumber    int32    `protobuf:"varint,4,opt,name=DeviceNumber,proto3" json:"DeviceNumber,omitempty"`
	UseExternalSNAT bool     `protobuf:"varint,5,opt,name=UseExternalSNAT,proto3" json:"UseExternalSNAT,omitempty"`
	VPCV4Cidrs      []string `protobuf:"bytes,6,rep,name=VPCV4cidrs,proto3" json:"VPCV4cidrs,omitempty"`
	VPCV6Cidrs      []string `protobuf:"bytes,7,rep,name=VPCV6cidrs,proto3" json:"VPCV6cidrs,omitempty"`
	// start of pod-eni parameters
	PodVlanId      int32  `protobuf:"varint,8,opt,name=PodVlanId,proto3" json:"PodVlanId,omitempty"`
	PodENIMAC      string `protobuf:"bytes,9,opt,name=PodENIMAC,proto3" json:"PodENIMAC,omitempty"`
	PodENISubnetGW string `protobuf:"bytes,10,opt,name=PodENISubnetGW,proto3" json:"PodENISubnetGW,omitempty"`
	ParentIfIndex  int32  `protobuf:"varint,11,opt,name=ParentIfIndex,proto3" json:"ParentIfIndex,omitempty"` // end of pod-eni parameters
}

func (x *AddNetworkReply) Reset() {
	*x = AddNetworkReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddNetworkReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddNetworkReply) ProtoMessage() {}

func (x *AddNetworkReply) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddNetworkReply.ProtoReflect.Descriptor instead.
func (*AddNetworkReply) Descriptor() ([]byte, []int) {
	return file_rpc_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *AddNetworkReply) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *AddNetworkReply) GetIPv4Addr() string {
	if x != nil {
		return x.IPv4Addr
	}
	return ""
}

func (x *AddNetworkReply) GetIPv6Addr() string {
	if x != nil {
		return x.IPv6Addr
	}
	return ""
}

func (x *AddNetworkReply) GetDeviceNumber() int32 {
	if x != nil {
		return x.DeviceNumber
	}
	return 0
}

func (x *AddNetworkReply) GetUseExternalSNAT() bool {
	if x != nil {
		return x.UseExternalSNAT
	}
	return false
}

func (x *AddNetworkReply) GetVPCV4Cidrs() []string {
	if x != nil {
		return x.VPCV4Cidrs
	}
	return nil
}

func (x *AddNetworkReply) GetVPCV6Cidrs() []string {
	if x != nil {
		return x.VPCV6Cidrs
	}
	return nil
}

func (x *AddNetworkReply) GetPodVlanId() int32 {
	if x != nil {
		return x.PodVlanId
	}
	return 0
}

func (x *AddNetworkReply) GetPodENIMAC() string {
	if x != nil {
		return x.PodENIMAC
	}
	return ""
}

func (x *AddNetworkReply) GetPodENISubnetGW() string {
	if x != nil {
		return x.PodENISubnetGW
	}
	return ""
}

func (x *AddNetworkReply) GetParentIfIndex() int32 {
	if x != nil {
		return x.ParentIfIndex
	}
	return 0
}

type DelNetworkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientVersion              string `protobuf:"bytes,9,opt,name=ClientVersion,proto3" json:"ClientVersion,omitempty"`
	K8S_POD_NAME               string `protobuf:"bytes,1,opt,name=K8S_POD_NAME,json=K8SPODNAME,proto3" json:"K8S_POD_NAME,omitempty"`
	K8S_POD_NAMESPACE          string `protobuf:"bytes,2,opt,name=K8S_POD_NAMESPACE,json=K8SPODNAMESPACE,proto3" json:"K8S_POD_NAMESPACE,omitempty"`
	K8S_POD_INFRA_CONTAINER_ID string `protobuf:"bytes,3,opt,name=K8S_POD_INFRA_CONTAINER_ID,json=K8SPODINFRACONTAINERID,proto3" json:"K8S_POD_INFRA_CONTAINER_ID,omitempty"`
	Reason                     string `protobuf:"bytes,5,opt,name=Reason,proto3" json:"Reason,omitempty"`
	ContainerID                string `protobuf:"bytes,8,opt,name=ContainerID,proto3" json:"ContainerID,omitempty"`
	IfName                     string `protobuf:"bytes,6,opt,name=IfName,proto3" json:"IfName,omitempty"`
	NetworkName                string `protobuf:"bytes,7,opt,name=NetworkName,proto3" json:"NetworkName,omitempty"` // next field: 10
}

func (x *DelNetworkRequest) Reset() {
	*x = DelNetworkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_rpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelNetworkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelNetworkRequest) ProtoMessage() {}

func (x *DelNetworkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_rpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelNetworkRequest.ProtoReflect.Descriptor instead.
func (*DelNetworkRequest) Descriptor() ([]byte, []int) {
	return file_rpc_rpc_proto_rawDescGZIP(), []int{2}
}

func (x *DelNetworkRequest) GetClientVersion() string {
	if x != nil {
		return x.ClientVersion
	}
	return ""
}

func (x *DelNetworkRequest) GetK8S_POD_NAME() string {
	if x != nil {
		return x.K8S_POD_NAME
	}
	return ""
}

func (x *DelNetworkRequest) GetK8S_POD_NAMESPACE() string {
	if x != nil {
		return x.K8S_POD_NAMESPACE
	}
	return ""
}

func (x *DelNetworkRequest) GetK8S_POD_INFRA_CONTAINER_ID() string {
	if x != nil {
		return x.K8S_POD_INFRA_CONTAINER_ID
	}
	return ""
}

func (x *DelNetworkRequest) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *DelNetworkRequest) GetContainerID() string {
	if x != nil {
		return x.ContainerID
	}
	return ""
}

func (x *DelNetworkRequest) GetIfName() string {
	if x != nil {
		return x.IfName
	}
	return ""
}

func (x *DelNetworkRequest) GetNetworkName() string {
	if x != nil {
		return x.NetworkName
	}
	return ""
}

type DelNetworkReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success      bool   `protobuf:"varint,1,opt,name=Success,proto3" json:"Success,omitempty"`
	IPv4Addr     string `protobuf:"bytes,2,opt,name=IPv4Addr,proto3" json:"IPv4Addr,omitempty"`
	DeviceNumber int32  `protobuf:"varint,3,opt,name=DeviceNumber,proto3" json:"DeviceNumber,omitempty"`
	// start of pod-eni parameters
	PodVlanId int32 `protobuf:"varint,4,opt,name=PodVlanId,proto3" json:"PodVlanId,omitempty"` // end of pod-eni parameters
}

func (x *DelNetworkReply) Reset() {
	*x = DelNetworkReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_rpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DelNetworkReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DelNetworkReply) ProtoMessage() {}

func (x *DelNetworkReply) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_rpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DelNetworkReply.ProtoReflect.Descriptor instead.
func (*DelNetworkReply) Descriptor() ([]byte, []int) {
	return file_rpc_rpc_proto_rawDescGZIP(), []int{3}
}

func (x *DelNetworkReply) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *DelNetworkReply) GetIPv4Addr() string {
	if x != nil {
		return x.IPv4Addr
	}
	return ""
}

func (x *DelNetworkReply) GetDeviceNumber() int32 {
	if x != nil {
		return x.DeviceNumber
	}
	return 0
}

func (x *DelNetworkReply) GetPodVlanId() int32 {
	if x != nil {
		return x.PodVlanId
	}
	return 0
}

var File_rpc_rpc_proto protoreflect.FileDescriptor

var file_rpc_rpc_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x72, 0x70, 0x63, 0x2f, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x72, 0x70, 0x63, 0x22, 0xb5, 0x02, 0x0a, 0x11, 0x41, 0x64, 0x64, 0x4e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x20, 0x0a, 0x0c, 0x4b, 0x38, 0x53, 0x5f, 0x50, 0x4f, 0x44, 0x5f, 0x4e, 0x41, 0x4d, 0x45,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x4b, 0x38, 0x53, 0x50, 0x4f, 0x44, 0x4e, 0x41,
	0x4d, 0x45, 0x12, 0x2a, 0x0a, 0x11, 0x4b, 0x38, 0x53, 0x5f, 0x50, 0x4f, 0x44, 0x5f, 0x4e, 0x41,
	0x4d, 0x45, 0x53, 0x50, 0x41, 0x43, 0x45, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x4b,
	0x38, 0x53, 0x50, 0x4f, 0x44, 0x4e, 0x41, 0x4d, 0x45, 0x53, 0x50, 0x41, 0x43, 0x45, 0x12, 0x3a,
	0x0a, 0x1a, 0x4b, 0x38, 0x53, 0x5f, 0x50, 0x4f, 0x44, 0x5f, 0x49, 0x4e, 0x46, 0x52, 0x41, 0x5f,
	0x43, 0x4f, 0x4e, 0x54, 0x41, 0x49, 0x4e, 0x45, 0x52, 0x5f, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x16, 0x4b, 0x38, 0x53, 0x50, 0x4f, 0x44, 0x49, 0x4e, 0x46, 0x52, 0x41, 0x43,
	0x4f, 0x4e, 0x54, 0x41, 0x49, 0x4e, 0x45, 0x52, 0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x6f,
	0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x49, 0x44, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06,
	0x49, 0x66, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x49, 0x66,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x4e, 0x65, 0x74, 0x6e, 0x73, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x4e, 0x65, 0x74, 0x6e, 0x73, 0x22, 0xfb, 0x02, 0x0a,
	0x0f, 0x41, 0x64, 0x64, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x49, 0x50,
	0x76, 0x34, 0x41, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x49, 0x50,
	0x76, 0x34, 0x41, 0x64, 0x64, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x49, 0x50, 0x76, 0x36, 0x41, 0x64,
	0x64, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x49, 0x50, 0x76, 0x36, 0x41, 0x64,
	0x64, 0x72, 0x12, 0x22, 0x0a, 0x0c, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x28, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x45, 0x78, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x53, 0x4e, 0x41, 0x54, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0f, 0x55, 0x73, 0x65, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x53, 0x4e, 0x41, 0x54,
	0x12, 0x1e, 0x0a, 0x0a, 0x56, 0x50, 0x43, 0x56, 0x34, 0x63, 0x69, 0x64, 0x72, 0x73, 0x18, 0x06,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x56, 0x50, 0x43, 0x56, 0x34, 0x63, 0x69, 0x64, 0x72, 0x73,
	0x12, 0x1e, 0x0a, 0x0a, 0x56, 0x50, 0x43, 0x56, 0x36, 0x63, 0x69, 0x64, 0x72, 0x73, 0x18, 0x07,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x56, 0x50, 0x43, 0x56, 0x36, 0x63, 0x69, 0x64, 0x72, 0x73,
	0x12, 0x1c, 0x0a, 0x09, 0x50, 0x6f, 0x64, 0x56, 0x6c, 0x61, 0x6e, 0x49, 0x64, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x09, 0x50, 0x6f, 0x64, 0x56, 0x6c, 0x61, 0x6e, 0x49, 0x64, 0x12, 0x1c,
	0x0a, 0x09, 0x50, 0x6f, 0x64, 0x45, 0x4e, 0x49, 0x4d, 0x41, 0x43, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x50, 0x6f, 0x64, 0x45, 0x4e, 0x49, 0x4d, 0x41, 0x43, 0x12, 0x26, 0x0a, 0x0e,
	0x50, 0x6f, 0x64, 0x45, 0x4e, 0x49, 0x53, 0x75, 0x62, 0x6e, 0x65, 0x74, 0x47, 0x57, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x50, 0x6f, 0x64, 0x45, 0x4e, 0x49, 0x53, 0x75, 0x62, 0x6e,
	0x65, 0x74, 0x47, 0x57, 0x12, 0x24, 0x0a, 0x0d, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x66,
	0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x50, 0x61, 0x72,
	0x65, 0x6e, 0x74, 0x49, 0x66, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x22, 0xb7, 0x02, 0x0a, 0x11, 0x44,
	0x65, 0x6c, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x24, 0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0c, 0x4b, 0x38, 0x53, 0x5f, 0x50, 0x4f,
	0x44, 0x5f, 0x4e, 0x41, 0x4d, 0x45, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x4b, 0x38,
	0x53, 0x50, 0x4f, 0x44, 0x4e, 0x41, 0x4d, 0x45, 0x12, 0x2a, 0x0a, 0x11, 0x4b, 0x38, 0x53, 0x5f,
	0x50, 0x4f, 0x44, 0x5f, 0x4e, 0x41, 0x4d, 0x45, 0x53, 0x50, 0x41, 0x43, 0x45, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x4b, 0x38, 0x53, 0x50, 0x4f, 0x44, 0x4e, 0x41, 0x4d, 0x45, 0x53,
	0x50, 0x41, 0x43, 0x45, 0x12, 0x3a, 0x0a, 0x1a, 0x4b, 0x38, 0x53, 0x5f, 0x50, 0x4f, 0x44, 0x5f,
	0x49, 0x4e, 0x46, 0x52, 0x41, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x41, 0x49, 0x4e, 0x45, 0x52, 0x5f,
	0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x16, 0x4b, 0x38, 0x53, 0x50, 0x4f, 0x44,
	0x49, 0x4e, 0x46, 0x52, 0x41, 0x43, 0x4f, 0x4e, 0x54, 0x41, 0x49, 0x4e, 0x45, 0x52, 0x49, 0x44,
	0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x6f, 0x6e, 0x74,
	0x61, 0x69, 0x6e, 0x65, 0x72, 0x49, 0x44, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x43,
	0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x49, 0x66,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x49, 0x66, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x4e, 0x61, 0x6d, 0x65, 0x22, 0x89, 0x01, 0x0a, 0x0f, 0x44, 0x65, 0x6c, 0x4e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x49, 0x50, 0x76, 0x34, 0x41, 0x64, 0x64, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x49, 0x50, 0x76, 0x34, 0x41, 0x64, 0x64, 0x72, 0x12, 0x22,
	0x0a, 0x0c, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x50, 0x6f, 0x64, 0x56, 0x6c, 0x61, 0x6e, 0x49, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x50, 0x6f, 0x64, 0x56, 0x6c, 0x61, 0x6e, 0x49, 0x64,
	0x32, 0x88, 0x01, 0x0a, 0x0a, 0x43, 0x4e, 0x49, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x12,
	0x3c, 0x0a, 0x0a, 0x41, 0x64, 0x64, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x16, 0x2e,
	0x72, 0x70, 0x63, 0x2e, 0x41, 0x64, 0x64, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x41, 0x64, 0x64, 0x4e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x3c, 0x0a,
	0x0a, 0x44, 0x65, 0x6c, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x16, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x44, 0x65, 0x6c, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x44, 0x65, 0x6c, 0x4e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x72,
	0x70, 0x63, 0x2f, 0x3b, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_rpc_proto_rawDescOnce sync.Once
	file_rpc_rpc_proto_rawDescData = file_rpc_rpc_proto_rawDesc
)

func file_rpc_rpc_proto_rawDescGZIP() []byte {
	file_rpc_rpc_proto_rawDescOnce.Do(func() {
		file_rpc_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_rpc_proto_rawDescData)
	})
	return file_rpc_rpc_proto_rawDescData
}

var file_rpc_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_rpc_rpc_proto_goTypes = []interface{}{
	(*AddNetworkRequest)(nil), // 0: rpc.AddNetworkRequest
	(*AddNetworkReply)(nil),   // 1: rpc.AddNetworkReply
	(*DelNetworkRequest)(nil), // 2: rpc.DelNetworkRequest
	(*DelNetworkReply)(nil),   // 3: rpc.DelNetworkReply
}
var file_rpc_rpc_proto_depIdxs = []int32{
	0, // 0: rpc.CNIBackend.AddNetwork:input_type -> rpc.AddNetworkRequest
	2, // 1: rpc.CNIBackend.DelNetwork:input_type -> rpc.DelNetworkRequest
	1, // 2: rpc.CNIBackend.AddNetwork:output_type -> rpc.AddNetworkReply
	3, // 3: rpc.CNIBackend.DelNetwork:output_type -> rpc.DelNetworkReply
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rpc_rpc_proto_init() }
func file_rpc_rpc_proto_init() {
	if File_rpc_rpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddNetworkRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddNetworkReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_rpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelNetworkRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_rpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DelNetworkReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_rpc_proto_goTypes,
		DependencyIndexes: file_rpc_rpc_proto_depIdxs,
		MessageInfos:      file_rpc_rpc_proto_msgTypes,
	}.Build()
	File_rpc_rpc_proto = out.File
	file_rpc_rpc_proto_rawDesc = nil
	file_rpc_rpc_proto_goTypes = nil
	file_rpc_rpc_proto_depIdxs = nil
}
