// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: peering.proto

package peering

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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_peering_proto_rawDescGZIP(), []int{0}
}

type CreateKeygroupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup string `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
	Expiry   int64  `protobuf:"varint,2,opt,name=expiry,proto3" json:"expiry,omitempty"`
}

func (x *CreateKeygroupRequest) Reset() {
	*x = CreateKeygroupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateKeygroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateKeygroupRequest) ProtoMessage() {}

func (x *CreateKeygroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateKeygroupRequest.ProtoReflect.Descriptor instead.
func (*CreateKeygroupRequest) Descriptor() ([]byte, []int) {
	return file_peering_proto_rawDescGZIP(), []int{1}
}

func (x *CreateKeygroupRequest) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

func (x *CreateKeygroupRequest) GetExpiry() int64 {
	if x != nil {
		return x.Expiry
	}
	return 0
}

type DeleteKeygroupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup string `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
}

func (x *DeleteKeygroupRequest) Reset() {
	*x = DeleteKeygroupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteKeygroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteKeygroupRequest) ProtoMessage() {}

func (x *DeleteKeygroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteKeygroupRequest.ProtoReflect.Descriptor instead.
func (*DeleteKeygroupRequest) Descriptor() ([]byte, []int) {
	return file_peering_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteKeygroupRequest) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

type PutItemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup   string            `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
	Id         string            `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Val        string            `protobuf:"bytes,3,opt,name=val,proto3" json:"val,omitempty"`
	Tombstoned bool              `protobuf:"varint,4,opt,name=tombstoned,proto3" json:"tombstoned,omitempty"`
	Version    map[string]uint64 `protobuf:"bytes,5,rep,name=version,proto3" json:"version,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *PutItemRequest) Reset() {
	*x = PutItemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutItemRequest) ProtoMessage() {}

func (x *PutItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutItemRequest.ProtoReflect.Descriptor instead.
func (*PutItemRequest) Descriptor() ([]byte, []int) {
	return file_peering_proto_rawDescGZIP(), []int{3}
}

func (x *PutItemRequest) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

func (x *PutItemRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PutItemRequest) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

func (x *PutItemRequest) GetTombstoned() bool {
	if x != nil {
		return x.Tombstoned
	}
	return false
}

func (x *PutItemRequest) GetVersion() map[string]uint64 {
	if x != nil {
		return x.Version
	}
	return nil
}

type GetItemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup string `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
	Id       string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetItemRequest) Reset() {
	*x = GetItemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetItemRequest) ProtoMessage() {}

func (x *GetItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetItemRequest.ProtoReflect.Descriptor instead.
func (*GetItemRequest) Descriptor() ([]byte, []int) {
	return file_peering_proto_rawDescGZIP(), []int{4}
}

func (x *GetItemRequest) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

func (x *GetItemRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetItemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*Data `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *GetItemResponse) Reset() {
	*x = GetItemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetItemResponse) ProtoMessage() {}

func (x *GetItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetItemResponse.ProtoReflect.Descriptor instead.
func (*GetItemResponse) Descriptor() ([]byte, []int) {
	return file_peering_proto_rawDescGZIP(), []int{5}
}

func (x *GetItemResponse) GetData() []*Data {
	if x != nil {
		return x.Data
	}
	return nil
}

type GetAllItemsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup string `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
}

func (x *GetAllItemsRequest) Reset() {
	*x = GetAllItemsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllItemsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllItemsRequest) ProtoMessage() {}

func (x *GetAllItemsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllItemsRequest.ProtoReflect.Descriptor instead.
func (*GetAllItemsRequest) Descriptor() ([]byte, []int) {
	return file_peering_proto_rawDescGZIP(), []int{6}
}

func (x *GetAllItemsRequest) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

type GetAllItemsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*Data `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *GetAllItemsResponse) Reset() {
	*x = GetAllItemsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllItemsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllItemsResponse) ProtoMessage() {}

func (x *GetAllItemsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllItemsResponse.ProtoReflect.Descriptor instead.
func (*GetAllItemsResponse) Descriptor() ([]byte, []int) {
	return file_peering_proto_rawDescGZIP(), []int{7}
}

func (x *GetAllItemsResponse) GetData() []*Data {
	if x != nil {
		return x.Data
	}
	return nil
}

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Val     string            `protobuf:"bytes,2,opt,name=val,proto3" json:"val,omitempty"`
	Version map[string]uint64 `protobuf:"bytes,3,rep,name=version,proto3" json:"version,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_peering_proto_rawDescGZIP(), []int{8}
}

func (x *Data) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Data) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

func (x *Data) GetVersion() map[string]uint64 {
	if x != nil {
		return x.Version
	}
	return nil
}

type UpdateItemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup string            `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
	Id       string            `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Val      string            `protobuf:"bytes,3,opt,name=val,proto3" json:"val,omitempty"`
	Version  map[string]uint64 `protobuf:"bytes,4,rep,name=version,proto3" json:"version,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *UpdateItemRequest) Reset() {
	*x = UpdateItemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateItemRequest) ProtoMessage() {}

func (x *UpdateItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateItemRequest.ProtoReflect.Descriptor instead.
func (*UpdateItemRequest) Descriptor() ([]byte, []int) {
	return file_peering_proto_rawDescGZIP(), []int{9}
}

func (x *UpdateItemRequest) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

func (x *UpdateItemRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateItemRequest) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

func (x *UpdateItemRequest) GetVersion() map[string]uint64 {
	if x != nil {
		return x.Version
	}
	return nil
}

type AppendItemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup string `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
	Id       string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Data     string `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *AppendItemRequest) Reset() {
	*x = AppendItemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendItemRequest) ProtoMessage() {}

func (x *AppendItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendItemRequest.ProtoReflect.Descriptor instead.
func (*AppendItemRequest) Descriptor() ([]byte, []int) {
	return file_peering_proto_rawDescGZIP(), []int{10}
}

func (x *AppendItemRequest) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

func (x *AppendItemRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AppendItemRequest) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

var File_peering_proto protoreflect.FileDescriptor

var file_peering_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x10, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e,
	0x67, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x4b, 0x0a, 0x15, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12,
	0x16, 0x0a, 0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x22, 0x33, 0x0a, 0x15, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x22, 0xf3, 0x01, 0x0a,
	0x0e, 0x50, 0x75, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x76,
	0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x12, 0x1e, 0x0a,
	0x0a, 0x74, 0x6f, 0x6d, 0x62, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0a, 0x74, 0x6f, 0x6d, 0x62, 0x73, 0x74, 0x6f, 0x6e, 0x65, 0x64, 0x12, 0x47, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2d,
	0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e,
	0x67, 0x2e, 0x50, 0x75, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x3a, 0x0a, 0x0c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x3c, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x22, 0x3d, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x65, 0x65,
	0x72, 0x69, 0x6e, 0x67, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x30, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x22, 0x41, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x49, 0x74, 0x65, 0x6d, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65,
	0x64, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0xa3, 0x01, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x12,
	0x3d, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x23, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x65, 0x65, 0x72,
	0x69, 0x6e, 0x67, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x3a,
	0x0a, 0x0c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xd9, 0x01, 0x0a, 0x11, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x76, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x12, 0x4a,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x30, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x69,
	0x6e, 0x67, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x3a, 0x0a, 0x0c, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x53, 0x0a, 0x11, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64,
	0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6b,
	0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b,
	0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0xec, 0x03, 0x0a, 0x04,
	0x4e, 0x6f, 0x64, 0x65, 0x12, 0x52, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4b, 0x65,
	0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x27, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65,
	0x64, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x17, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x69,
	0x6e, 0x67, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x52, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x27, 0x2e, 0x6d, 0x63, 0x63,
	0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x70,
	0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x44, 0x0a, 0x07,
	0x50, 0x75, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x20, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72,
	0x65, 0x64, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x75, 0x74, 0x49, 0x74,
	0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6d, 0x63, 0x63, 0x2e,
	0x66, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x4a, 0x0a, 0x0a, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x49, 0x74, 0x65, 0x6d,
	0x12, 0x23, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x65, 0x65, 0x72,
	0x69, 0x6e, 0x67, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64,
	0x2e, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x4e,
	0x0a, 0x07, 0x47, 0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x20, 0x2e, 0x6d, 0x63, 0x63, 0x2e,
	0x66, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x47, 0x65, 0x74,
	0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x6d, 0x63,
	0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x47,
	0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5a,
	0x0a, 0x0b, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x24, 0x2e,
	0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67,
	0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x70,
	0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x49, 0x74, 0x65,
	0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x3b,
	0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_peering_proto_rawDescOnce sync.Once
	file_peering_proto_rawDescData = file_peering_proto_rawDesc
)

func file_peering_proto_rawDescGZIP() []byte {
	file_peering_proto_rawDescOnce.Do(func() {
		file_peering_proto_rawDescData = protoimpl.X.CompressGZIP(file_peering_proto_rawDescData)
	})
	return file_peering_proto_rawDescData
}

var file_peering_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_peering_proto_goTypes = []interface{}{
	(*Empty)(nil),                 // 0: mcc.fred.peering.Empty
	(*CreateKeygroupRequest)(nil), // 1: mcc.fred.peering.CreateKeygroupRequest
	(*DeleteKeygroupRequest)(nil), // 2: mcc.fred.peering.DeleteKeygroupRequest
	(*PutItemRequest)(nil),        // 3: mcc.fred.peering.PutItemRequest
	(*GetItemRequest)(nil),        // 4: mcc.fred.peering.GetItemRequest
	(*GetItemResponse)(nil),       // 5: mcc.fred.peering.GetItemResponse
	(*GetAllItemsRequest)(nil),    // 6: mcc.fred.peering.GetAllItemsRequest
	(*GetAllItemsResponse)(nil),   // 7: mcc.fred.peering.GetAllItemsResponse
	(*Data)(nil),                  // 8: mcc.fred.peering.Data
	(*UpdateItemRequest)(nil),     // 9: mcc.fred.peering.UpdateItemRequest
	(*AppendItemRequest)(nil),     // 10: mcc.fred.peering.AppendItemRequest
	nil,                           // 11: mcc.fred.peering.PutItemRequest.VersionEntry
	nil,                           // 12: mcc.fred.peering.Data.VersionEntry
	nil,                           // 13: mcc.fred.peering.UpdateItemRequest.VersionEntry
}
var file_peering_proto_depIdxs = []int32{
	11, // 0: mcc.fred.peering.PutItemRequest.version:type_name -> mcc.fred.peering.PutItemRequest.VersionEntry
	8,  // 1: mcc.fred.peering.GetItemResponse.data:type_name -> mcc.fred.peering.Data
	8,  // 2: mcc.fred.peering.GetAllItemsResponse.data:type_name -> mcc.fred.peering.Data
	12, // 3: mcc.fred.peering.Data.version:type_name -> mcc.fred.peering.Data.VersionEntry
	13, // 4: mcc.fred.peering.UpdateItemRequest.version:type_name -> mcc.fred.peering.UpdateItemRequest.VersionEntry
	1,  // 5: mcc.fred.peering.Node.CreateKeygroup:input_type -> mcc.fred.peering.CreateKeygroupRequest
	2,  // 6: mcc.fred.peering.Node.DeleteKeygroup:input_type -> mcc.fred.peering.DeleteKeygroupRequest
	3,  // 7: mcc.fred.peering.Node.PutItem:input_type -> mcc.fred.peering.PutItemRequest
	10, // 8: mcc.fred.peering.Node.AppendItem:input_type -> mcc.fred.peering.AppendItemRequest
	4,  // 9: mcc.fred.peering.Node.GetItem:input_type -> mcc.fred.peering.GetItemRequest
	6,  // 10: mcc.fred.peering.Node.GetAllItems:input_type -> mcc.fred.peering.GetAllItemsRequest
	0,  // 11: mcc.fred.peering.Node.CreateKeygroup:output_type -> mcc.fred.peering.Empty
	0,  // 12: mcc.fred.peering.Node.DeleteKeygroup:output_type -> mcc.fred.peering.Empty
	0,  // 13: mcc.fred.peering.Node.PutItem:output_type -> mcc.fred.peering.Empty
	0,  // 14: mcc.fred.peering.Node.AppendItem:output_type -> mcc.fred.peering.Empty
	5,  // 15: mcc.fred.peering.Node.GetItem:output_type -> mcc.fred.peering.GetItemResponse
	7,  // 16: mcc.fred.peering.Node.GetAllItems:output_type -> mcc.fred.peering.GetAllItemsResponse
	11, // [11:17] is the sub-list for method output_type
	5,  // [5:11] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_peering_proto_init() }
func file_peering_proto_init() {
	if File_peering_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_peering_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_peering_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateKeygroupRequest); i {
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
		file_peering_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteKeygroupRequest); i {
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
		file_peering_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutItemRequest); i {
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
		file_peering_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetItemRequest); i {
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
		file_peering_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetItemResponse); i {
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
		file_peering_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllItemsRequest); i {
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
		file_peering_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllItemsResponse); i {
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
		file_peering_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data); i {
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
		file_peering_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateItemRequest); i {
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
		file_peering_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendItemRequest); i {
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
			RawDescriptor: file_peering_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_peering_proto_goTypes,
		DependencyIndexes: file_peering_proto_depIdxs,
		MessageInfos:      file_peering_proto_msgTypes,
	}.Build()
	File_peering_proto = out.File
	file_peering_proto_rawDesc = nil
	file_peering_proto_goTypes = nil
	file_peering_proto_depIdxs = nil
}
