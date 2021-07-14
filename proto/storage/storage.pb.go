// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: storage.proto

package storage

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

type Item struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup string `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
	Id       string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Val      string `protobuf:"bytes,3,opt,name=val,proto3" json:"val,omitempty"`
}

func (x *Item) Reset() {
	*x = Item{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Item) ProtoMessage() {}

func (x *Item) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Item.ProtoReflect.Descriptor instead.
func (*Item) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{0}
}

func (x *Item) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

func (x *Item) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Item) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

type ScanRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   *Key   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Count uint64 `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *ScanRequest) Reset() {
	*x = ScanRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScanRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanRequest) ProtoMessage() {}

func (x *ScanRequest) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanRequest.ProtoReflect.Descriptor instead.
func (*ScanRequest) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{1}
}

func (x *ScanRequest) GetKey() *Key {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *ScanRequest) GetCount() uint64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type UpdateItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup string `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
	Id       string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Val      string `protobuf:"bytes,3,opt,name=val,proto3" json:"val,omitempty"`
	Expiry   int64  `protobuf:"varint,4,opt,name=expiry,proto3" json:"expiry,omitempty"`
}

func (x *UpdateItem) Reset() {
	*x = UpdateItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateItem) ProtoMessage() {}

func (x *UpdateItem) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateItem.ProtoReflect.Descriptor instead.
func (*UpdateItem) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateItem) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

func (x *UpdateItem) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateItem) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

func (x *UpdateItem) GetExpiry() int64 {
	if x != nil {
		return x.Expiry
	}
	return 0
}

type AppendItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup string `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
	Val      string `protobuf:"bytes,2,opt,name=val,proto3" json:"val,omitempty"`
	Expiry   int64  `protobuf:"varint,3,opt,name=expiry,proto3" json:"expiry,omitempty"`
}

func (x *AppendItem) Reset() {
	*x = AppendItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendItem) ProtoMessage() {}

func (x *AppendItem) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendItem.ProtoReflect.Descriptor instead.
func (*AppendItem) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{3}
}

func (x *AppendItem) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

func (x *AppendItem) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

func (x *AppendItem) GetExpiry() int64 {
	if x != nil {
		return x.Expiry
	}
	return 0
}

type Trigger struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Host string `protobuf:"bytes,2,opt,name=host,proto3" json:"host,omitempty"`
}

func (x *Trigger) Reset() {
	*x = Trigger{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Trigger) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Trigger) ProtoMessage() {}

func (x *Trigger) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Trigger.ProtoReflect.Descriptor instead.
func (*Trigger) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{4}
}

func (x *Trigger) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Trigger) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

// A Key uniquely identifies data. In our case it contains a keygroup and the id
type Key struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup string `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
	Id       string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Key) Reset() {
	*x = Key{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Key) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Key) ProtoMessage() {}

func (x *Key) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Key.ProtoReflect.Descriptor instead.
func (*Key) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{5}
}

func (x *Key) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

func (x *Key) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Val struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Val string `protobuf:"bytes,1,opt,name=val,proto3" json:"val,omitempty"`
}

func (x *Val) Reset() {
	*x = Val{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Val) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Val) ProtoMessage() {}

func (x *Val) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Val.ProtoReflect.Descriptor instead.
func (*Val) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{6}
}

func (x *Val) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

type Keygroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup string `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
}

func (x *Keygroup) Reset() {
	*x = Keygroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Keygroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Keygroup) ProtoMessage() {}

func (x *Keygroup) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Keygroup.ProtoReflect.Descriptor instead.
func (*Keygroup) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{7}
}

func (x *Keygroup) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

type KeygroupTrigger struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup string   `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
	Trigger  *Trigger `protobuf:"bytes,2,opt,name=trigger,proto3" json:"trigger,omitempty"`
}

func (x *KeygroupTrigger) Reset() {
	*x = KeygroupTrigger{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeygroupTrigger) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeygroupTrigger) ProtoMessage() {}

func (x *KeygroupTrigger) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeygroupTrigger.ProtoReflect.Descriptor instead.
func (*KeygroupTrigger) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{8}
}

func (x *KeygroupTrigger) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

func (x *KeygroupTrigger) GetTrigger() *Trigger {
	if x != nil {
		return x.Trigger
	}
	return nil
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{9}
}

func (x *Response) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *Response) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_storage_proto protoreflect.FileDescriptor

var file_storage_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x10, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x22, 0x44, 0x0a, 0x04, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x22, 0x4c, 0x0a, 0x0b, 0x53, 0x63, 0x61, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x4b, 0x65, 0x79, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x62, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x49,
	0x74, 0x65, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x76, 0x61,
	0x6c, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x22, 0x52, 0x0a, 0x0a, 0x41, 0x70, 0x70,
	0x65, 0x6e, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x76, 0x61, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x22, 0x2d, 0x0a,
	0x07, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x22, 0x31, 0x0a, 0x03,
	0x4b, 0x65, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x17, 0x0a, 0x03, 0x56, 0x61, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x22, 0x26, 0x0a, 0x08, 0x4b, 0x65, 0x79, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x22, 0x62, 0x0a, 0x0f, 0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x72, 0x69, 0x67,
	0x67, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12,
	0x33, 0x0a, 0x07, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x2e, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x52, 0x07, 0x74, 0x72, 0x69,
	0x67, 0x67, 0x65, 0x72, 0x22, 0x3e, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x32, 0xf1, 0x07, 0x0a, 0x08, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73,
	0x65, 0x12, 0x44, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1c, 0x2e, 0x6d, 0x63,
	0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x1a, 0x1a, 0x2e, 0x6d, 0x63, 0x63, 0x2e,
	0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3d, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x12, 0x15, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x2e, 0x4b, 0x65, 0x79, 0x1a, 0x1a, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66,
	0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x06, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64,
	0x12, 0x1c, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x1a, 0x15,
	0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x2e, 0x4b, 0x65, 0x79, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x04, 0x52, 0x65, 0x61, 0x64, 0x12,
	0x15, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x2e, 0x4b, 0x65, 0x79, 0x1a, 0x15, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65,
	0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x56, 0x61, 0x6c, 0x22, 0x00, 0x12,
	0x41, 0x0a, 0x04, 0x53, 0x63, 0x61, 0x6e, 0x12, 0x1d, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72,
	0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x53, 0x63, 0x61, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65,
	0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x22, 0x00,
	0x30, 0x01, 0x12, 0x41, 0x0a, 0x07, 0x52, 0x65, 0x61, 0x64, 0x41, 0x6c, 0x6c, 0x12, 0x1a, 0x2e,
	0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x2e, 0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x1a, 0x16, 0x2e, 0x6d, 0x63, 0x63, 0x2e,
	0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x49, 0x74, 0x65,
	0x6d, 0x22, 0x00, 0x30, 0x01, 0x12, 0x3c, 0x0a, 0x03, 0x49, 0x44, 0x73, 0x12, 0x1a, 0x2e, 0x6d,
	0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e,
	0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x1a, 0x15, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66,
	0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x4b, 0x65, 0x79, 0x22,
	0x00, 0x30, 0x01, 0x12, 0x3d, 0x0a, 0x06, 0x45, 0x78, 0x69, 0x73, 0x74, 0x73, 0x12, 0x15, 0x2e,
	0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x2e, 0x4b, 0x65, 0x79, 0x1a, 0x1a, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e,
	0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x4a, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x12, 0x1a, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e,
	0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x1a, 0x1a, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4a,
	0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x12, 0x1a, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x2e, 0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x1a, 0x1a, 0x2e, 0x6d,
	0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4a, 0x0a, 0x0e, 0x45, 0x78,
	0x69, 0x73, 0x74, 0x73, 0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x1a, 0x2e, 0x6d,
	0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e,
	0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x1a, 0x1a, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66,
	0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x4b, 0x65, 0x79,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x12, 0x21, 0x2e, 0x6d,
	0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e,
	0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x1a,
	0x1a, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x58, 0x0a,
	0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x54,
	0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x12, 0x21, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65,
	0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x1a, 0x1a, 0x2e, 0x6d, 0x63, 0x63, 0x2e,
	0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4f, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x4b, 0x65,
	0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x12, 0x1a, 0x2e,
	0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x2e, 0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x1a, 0x19, 0x2e, 0x6d, 0x63, 0x63, 0x2e,
	0x66, 0x72, 0x65, 0x64, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x54, 0x72, 0x69,
	0x67, 0x67, 0x65, 0x72, 0x22, 0x00, 0x30, 0x01, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x3b, 0x73, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_storage_proto_rawDescOnce sync.Once
	file_storage_proto_rawDescData = file_storage_proto_rawDesc
)

func file_storage_proto_rawDescGZIP() []byte {
	file_storage_proto_rawDescOnce.Do(func() {
		file_storage_proto_rawDescData = protoimpl.X.CompressGZIP(file_storage_proto_rawDescData)
	})
	return file_storage_proto_rawDescData
}

var file_storage_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_storage_proto_goTypes = []interface{}{
	(*Item)(nil),            // 0: mcc.fred.storage.Item
	(*ScanRequest)(nil),     // 1: mcc.fred.storage.ScanRequest
	(*UpdateItem)(nil),      // 2: mcc.fred.storage.UpdateItem
	(*AppendItem)(nil),      // 3: mcc.fred.storage.AppendItem
	(*Trigger)(nil),         // 4: mcc.fred.storage.Trigger
	(*Key)(nil),             // 5: mcc.fred.storage.Key
	(*Val)(nil),             // 6: mcc.fred.storage.Val
	(*Keygroup)(nil),        // 7: mcc.fred.storage.Keygroup
	(*KeygroupTrigger)(nil), // 8: mcc.fred.storage.KeygroupTrigger
	(*Response)(nil),        // 9: mcc.fred.storage.Response
}
var file_storage_proto_depIdxs = []int32{
	5,  // 0: mcc.fred.storage.ScanRequest.key:type_name -> mcc.fred.storage.Key
	4,  // 1: mcc.fred.storage.KeygroupTrigger.trigger:type_name -> mcc.fred.storage.Trigger
	2,  // 2: mcc.fred.storage.Database.Update:input_type -> mcc.fred.storage.UpdateItem
	5,  // 3: mcc.fred.storage.Database.Delete:input_type -> mcc.fred.storage.Key
	3,  // 4: mcc.fred.storage.Database.Append:input_type -> mcc.fred.storage.AppendItem
	5,  // 5: mcc.fred.storage.Database.Read:input_type -> mcc.fred.storage.Key
	1,  // 6: mcc.fred.storage.Database.Scan:input_type -> mcc.fred.storage.ScanRequest
	7,  // 7: mcc.fred.storage.Database.ReadAll:input_type -> mcc.fred.storage.Keygroup
	7,  // 8: mcc.fred.storage.Database.IDs:input_type -> mcc.fred.storage.Keygroup
	5,  // 9: mcc.fred.storage.Database.Exists:input_type -> mcc.fred.storage.Key
	7,  // 10: mcc.fred.storage.Database.CreateKeygroup:input_type -> mcc.fred.storage.Keygroup
	7,  // 11: mcc.fred.storage.Database.DeleteKeygroup:input_type -> mcc.fred.storage.Keygroup
	7,  // 12: mcc.fred.storage.Database.ExistsKeygroup:input_type -> mcc.fred.storage.Keygroup
	8,  // 13: mcc.fred.storage.Database.AddKeygroupTrigger:input_type -> mcc.fred.storage.KeygroupTrigger
	8,  // 14: mcc.fred.storage.Database.DeleteKeygroupTrigger:input_type -> mcc.fred.storage.KeygroupTrigger
	7,  // 15: mcc.fred.storage.Database.GetKeygroupTrigger:input_type -> mcc.fred.storage.Keygroup
	9,  // 16: mcc.fred.storage.Database.Update:output_type -> mcc.fred.storage.Response
	9,  // 17: mcc.fred.storage.Database.Delete:output_type -> mcc.fred.storage.Response
	5,  // 18: mcc.fred.storage.Database.Append:output_type -> mcc.fred.storage.Key
	6,  // 19: mcc.fred.storage.Database.Read:output_type -> mcc.fred.storage.Val
	0,  // 20: mcc.fred.storage.Database.Scan:output_type -> mcc.fred.storage.Item
	0,  // 21: mcc.fred.storage.Database.ReadAll:output_type -> mcc.fred.storage.Item
	5,  // 22: mcc.fred.storage.Database.IDs:output_type -> mcc.fred.storage.Key
	9,  // 23: mcc.fred.storage.Database.Exists:output_type -> mcc.fred.storage.Response
	9,  // 24: mcc.fred.storage.Database.CreateKeygroup:output_type -> mcc.fred.storage.Response
	9,  // 25: mcc.fred.storage.Database.DeleteKeygroup:output_type -> mcc.fred.storage.Response
	9,  // 26: mcc.fred.storage.Database.ExistsKeygroup:output_type -> mcc.fred.storage.Response
	9,  // 27: mcc.fred.storage.Database.AddKeygroupTrigger:output_type -> mcc.fred.storage.Response
	9,  // 28: mcc.fred.storage.Database.DeleteKeygroupTrigger:output_type -> mcc.fred.storage.Response
	4,  // 29: mcc.fred.storage.Database.GetKeygroupTrigger:output_type -> mcc.fred.storage.Trigger
	16, // [16:30] is the sub-list for method output_type
	2,  // [2:16] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_storage_proto_init() }
func file_storage_proto_init() {
	if File_storage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_storage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Item); i {
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
		file_storage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScanRequest); i {
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
		file_storage_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateItem); i {
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
		file_storage_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendItem); i {
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
		file_storage_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Trigger); i {
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
		file_storage_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Key); i {
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
		file_storage_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Val); i {
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
		file_storage_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Keygroup); i {
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
		file_storage_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeygroupTrigger); i {
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
		file_storage_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
			RawDescriptor: file_storage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_storage_proto_goTypes,
		DependencyIndexes: file_storage_proto_depIdxs,
		MessageInfos:      file_storage_proto_msgTypes,
	}.Build()
	File_storage_proto = out.File
	file_storage_proto_rawDesc = nil
	file_storage_proto_goTypes = nil
	file_storage_proto_depIdxs = nil
}
