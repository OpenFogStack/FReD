// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: trigger.proto

package trigger

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
		mi := &file_trigger_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_trigger_proto_msgTypes[0]
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
	return file_trigger_proto_rawDescGZIP(), []int{0}
}

type PutItemTriggerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup string `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
	Id       string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Val      string `protobuf:"bytes,3,opt,name=val,proto3" json:"val,omitempty"`
}

func (x *PutItemTriggerRequest) Reset() {
	*x = PutItemTriggerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trigger_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutItemTriggerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutItemTriggerRequest) ProtoMessage() {}

func (x *PutItemTriggerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trigger_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutItemTriggerRequest.ProtoReflect.Descriptor instead.
func (*PutItemTriggerRequest) Descriptor() ([]byte, []int) {
	return file_trigger_proto_rawDescGZIP(), []int{1}
}

func (x *PutItemTriggerRequest) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

func (x *PutItemTriggerRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PutItemTriggerRequest) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

type DeleteItemTriggerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keygroup string `protobuf:"bytes,1,opt,name=keygroup,proto3" json:"keygroup,omitempty"`
	Id       string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteItemTriggerRequest) Reset() {
	*x = DeleteItemTriggerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trigger_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteItemTriggerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteItemTriggerRequest) ProtoMessage() {}

func (x *DeleteItemTriggerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trigger_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteItemTriggerRequest.ProtoReflect.Descriptor instead.
func (*DeleteItemTriggerRequest) Descriptor() ([]byte, []int) {
	return file_trigger_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteItemTriggerRequest) GetKeygroup() string {
	if x != nil {
		return x.Keygroup
	}
	return ""
}

func (x *DeleteItemTriggerRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_trigger_proto protoreflect.FileDescriptor

var file_trigger_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x10, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65,
	0x72, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x55, 0x0a, 0x15, 0x50, 0x75,
	0x74, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x76, 0x61,
	0x6c, 0x22, 0x46, 0x0a, 0x18, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x54,
	0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a,
	0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x32, 0xbb, 0x01, 0x0a, 0x0b, 0x54, 0x72,
	0x69, 0x67, 0x67, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x52, 0x0a, 0x0e, 0x50, 0x75, 0x74,
	0x49, 0x74, 0x65, 0x6d, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x12, 0x27, 0x2e, 0x6d, 0x63,
	0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x50,
	0x75, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e,
	0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x58, 0x0a,
	0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x72, 0x69, 0x67, 0x67,
	0x65, 0x72, 0x12, 0x2a, 0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x74, 0x72,
	0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x49, 0x74, 0x65, 0x6d,
	0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17,
	0x2e, 0x6d, 0x63, 0x63, 0x2e, 0x66, 0x72, 0x65, 0x64, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65,
	0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x0b, 0x5a, 0x09, 0x2e, 0x3b, 0x74, 0x72, 0x69,
	0x67, 0x67, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_trigger_proto_rawDescOnce sync.Once
	file_trigger_proto_rawDescData = file_trigger_proto_rawDesc
)

func file_trigger_proto_rawDescGZIP() []byte {
	file_trigger_proto_rawDescOnce.Do(func() {
		file_trigger_proto_rawDescData = protoimpl.X.CompressGZIP(file_trigger_proto_rawDescData)
	})
	return file_trigger_proto_rawDescData
}

var file_trigger_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_trigger_proto_goTypes = []interface{}{
	(*Empty)(nil),                    // 0: mcc.fred.trigger.Empty
	(*PutItemTriggerRequest)(nil),    // 1: mcc.fred.trigger.PutItemTriggerRequest
	(*DeleteItemTriggerRequest)(nil), // 2: mcc.fred.trigger.DeleteItemTriggerRequest
}
var file_trigger_proto_depIdxs = []int32{
	1, // 0: mcc.fred.trigger.TriggerNode.PutItemTrigger:input_type -> mcc.fred.trigger.PutItemTriggerRequest
	2, // 1: mcc.fred.trigger.TriggerNode.DeleteItemTrigger:input_type -> mcc.fred.trigger.DeleteItemTriggerRequest
	0, // 2: mcc.fred.trigger.TriggerNode.PutItemTrigger:output_type -> mcc.fred.trigger.Empty
	0, // 3: mcc.fred.trigger.TriggerNode.DeleteItemTrigger:output_type -> mcc.fred.trigger.Empty
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_trigger_proto_init() }
func file_trigger_proto_init() {
	if File_trigger_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_trigger_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_trigger_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutItemTriggerRequest); i {
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
		file_trigger_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteItemTriggerRequest); i {
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
			RawDescriptor: file_trigger_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_trigger_proto_goTypes,
		DependencyIndexes: file_trigger_proto_depIdxs,
		MessageInfos:      file_trigger_proto_msgTypes,
	}.Build()
	File_trigger_proto = out.File
	file_trigger_proto_rawDesc = nil
	file_trigger_proto_goTypes = nil
	file_trigger_proto_depIdxs = nil
}
