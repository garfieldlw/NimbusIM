// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: pusher/pusher.proto

package pusher

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

type DelayEnum int32

const (
	DelayEnum__DelayEnumUnknown DelayEnum = 0
	DelayEnum_DelayEnum1m       DelayEnum = 1
	DelayEnum_DelayEnum3m       DelayEnum = 2
	DelayEnum_DelayEnum10m      DelayEnum = 3
	DelayEnum_DelayEnum15m      DelayEnum = 4
	DelayEnum_DelayEnum30m      DelayEnum = 5
)

// Enum value maps for DelayEnum.
var (
	DelayEnum_name = map[int32]string{
		0: "_DelayEnumUnknown",
		1: "DelayEnum1m",
		2: "DelayEnum3m",
		3: "DelayEnum10m",
		4: "DelayEnum15m",
		5: "DelayEnum30m",
	}
	DelayEnum_value = map[string]int32{
		"_DelayEnumUnknown": 0,
		"DelayEnum1m":       1,
		"DelayEnum3m":       2,
		"DelayEnum10m":      3,
		"DelayEnum15m":      4,
		"DelayEnum30m":      5,
	}
)

func (x DelayEnum) Enum() *DelayEnum {
	p := new(DelayEnum)
	*p = x
	return p
}

func (x DelayEnum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DelayEnum) Descriptor() protoreflect.EnumDescriptor {
	return file_pusher_pusher_proto_enumTypes[0].Descriptor()
}

func (DelayEnum) Type() protoreflect.EnumType {
	return &file_pusher_pusher_proto_enumTypes[0]
}

func (x DelayEnum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DelayEnum.Descriptor instead.
func (DelayEnum) EnumDescriptor() ([]byte, []int) {
	return file_pusher_pusher_proto_rawDescGZIP(), []int{0}
}

type PusherKafkaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Topic string `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	Data  string `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *PusherKafkaRequest) Reset() {
	*x = PusherKafkaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pusher_pusher_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PusherKafkaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PusherKafkaRequest) ProtoMessage() {}

func (x *PusherKafkaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pusher_pusher_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PusherKafkaRequest.ProtoReflect.Descriptor instead.
func (*PusherKafkaRequest) Descriptor() ([]byte, []int) {
	return file_pusher_pusher_proto_rawDescGZIP(), []int{0}
}

func (x *PusherKafkaRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PusherKafkaRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *PusherKafkaRequest) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type PusherKafkaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *PusherKafkaResponse) Reset() {
	*x = PusherKafkaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pusher_pusher_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PusherKafkaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PusherKafkaResponse) ProtoMessage() {}

func (x *PusherKafkaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pusher_pusher_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PusherKafkaResponse.ProtoReflect.Descriptor instead.
func (*PusherKafkaResponse) Descriptor() ([]byte, []int) {
	return file_pusher_pusher_proto_rawDescGZIP(), []int{1}
}

func (x *PusherKafkaResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type PusherKafkaDelayRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int64     `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Topic string    `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	Data  string    `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Delay DelayEnum `protobuf:"varint,4,opt,name=delay,proto3,enum=com.NimbusIM.proto.pusher.DelayEnum" json:"delay,omitempty"`
}

func (x *PusherKafkaDelayRequest) Reset() {
	*x = PusherKafkaDelayRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pusher_pusher_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PusherKafkaDelayRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PusherKafkaDelayRequest) ProtoMessage() {}

func (x *PusherKafkaDelayRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pusher_pusher_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PusherKafkaDelayRequest.ProtoReflect.Descriptor instead.
func (*PusherKafkaDelayRequest) Descriptor() ([]byte, []int) {
	return file_pusher_pusher_proto_rawDescGZIP(), []int{2}
}

func (x *PusherKafkaDelayRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PusherKafkaDelayRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *PusherKafkaDelayRequest) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

func (x *PusherKafkaDelayRequest) GetDelay() DelayEnum {
	if x != nil {
		return x.Delay
	}
	return DelayEnum__DelayEnumUnknown
}

type PusherKafkaDelayResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *PusherKafkaDelayResponse) Reset() {
	*x = PusherKafkaDelayResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pusher_pusher_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PusherKafkaDelayResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PusherKafkaDelayResponse) ProtoMessage() {}

func (x *PusherKafkaDelayResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pusher_pusher_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PusherKafkaDelayResponse.ProtoReflect.Descriptor instead.
func (*PusherKafkaDelayResponse) Descriptor() ([]byte, []int) {
	return file_pusher_pusher_proto_rawDescGZIP(), []int{3}
}

func (x *PusherKafkaDelayResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

var File_pusher_pusher_proto protoreflect.FileDescriptor

var file_pusher_pusher_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x75, 0x73, 0x68, 0x65, 0x72, 0x2f, 0x70, 0x75, 0x73, 0x68, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x6d, 0x62, 0x75,
	0x73, 0x49, 0x4d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x75, 0x73, 0x68, 0x65, 0x72,
	0x22, 0x4e, 0x0a, 0x12, 0x50, 0x75, 0x73, 0x68, 0x65, 0x72, 0x4b, 0x61, 0x66, 0x6b, 0x61, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x22, 0x25, 0x0a, 0x13, 0x50, 0x75, 0x73, 0x68, 0x65, 0x72, 0x4b, 0x61, 0x66, 0x6b, 0x61, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x8f, 0x01, 0x0a, 0x17, 0x50, 0x75, 0x73, 0x68,
	0x65, 0x72, 0x4b, 0x61, 0x66, 0x6b, 0x61, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x3a, 0x0a,
	0x05, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x24, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x6d, 0x62, 0x75, 0x73, 0x49, 0x4d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x70, 0x75, 0x73, 0x68, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x45, 0x6e,
	0x75, 0x6d, 0x52, 0x05, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x22, 0x2a, 0x0a, 0x18, 0x50, 0x75, 0x73,
	0x68, 0x65, 0x72, 0x4b, 0x61, 0x66, 0x6b, 0x61, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x2a, 0x7a, 0x0a, 0x09, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x45, 0x6e,
	0x75, 0x6d, 0x12, 0x15, 0x0a, 0x11, 0x5f, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x45, 0x6e, 0x75, 0x6d,
	0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x44, 0x65, 0x6c,
	0x61, 0x79, 0x45, 0x6e, 0x75, 0x6d, 0x31, 0x6d, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x44, 0x65,
	0x6c, 0x61, 0x79, 0x45, 0x6e, 0x75, 0x6d, 0x33, 0x6d, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x44,
	0x65, 0x6c, 0x61, 0x79, 0x45, 0x6e, 0x75, 0x6d, 0x31, 0x30, 0x6d, 0x10, 0x03, 0x12, 0x10, 0x0a,
	0x0c, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x45, 0x6e, 0x75, 0x6d, 0x31, 0x35, 0x6d, 0x10, 0x04, 0x12,
	0x10, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x45, 0x6e, 0x75, 0x6d, 0x33, 0x30, 0x6d, 0x10,
	0x05, 0x32, 0xfc, 0x01, 0x0a, 0x0b, 0x4b, 0x61, 0x66, 0x6b, 0x61, 0x50, 0x75, 0x73, 0x68, 0x65,
	0x72, 0x12, 0x6e, 0x0a, 0x0b, 0x50, 0x75, 0x73, 0x68, 0x54, 0x6f, 0x4b, 0x61, 0x66, 0x6b, 0x61,
	0x12, 0x2d, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x6d, 0x62, 0x75, 0x73, 0x49, 0x4d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x75, 0x73, 0x68, 0x65, 0x72, 0x2e, 0x50, 0x75, 0x73,
	0x68, 0x65, 0x72, 0x4b, 0x61, 0x66, 0x6b, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2e, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x6d, 0x62, 0x75, 0x73, 0x49, 0x4d, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x75, 0x73, 0x68, 0x65, 0x72, 0x2e, 0x50, 0x75, 0x73, 0x68,
	0x65, 0x72, 0x4b, 0x61, 0x66, 0x6b, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x7d, 0x0a, 0x10, 0x50, 0x75, 0x73, 0x68, 0x54, 0x6f, 0x4b, 0x61, 0x66, 0x6b, 0x61,
	0x44, 0x65, 0x6c, 0x61, 0x79, 0x12, 0x32, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x6d, 0x62,
	0x75, 0x73, 0x49, 0x4d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x75, 0x73, 0x68, 0x65,
	0x72, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x65, 0x72, 0x4b, 0x61, 0x66, 0x6b, 0x61, 0x44, 0x65, 0x6c,
	0x61, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x63, 0x6f, 0x6d, 0x2e,
	0x4e, 0x69, 0x6d, 0x62, 0x75, 0x73, 0x49, 0x4d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70,
	0x75, 0x73, 0x68, 0x65, 0x72, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x65, 0x72, 0x4b, 0x61, 0x66, 0x6b,
	0x61, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x09, 0x5a, 0x07, 0x70, 0x75, 0x73, 0x68, 0x65, 0x72, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_pusher_pusher_proto_rawDescOnce sync.Once
	file_pusher_pusher_proto_rawDescData = file_pusher_pusher_proto_rawDesc
)

func file_pusher_pusher_proto_rawDescGZIP() []byte {
	file_pusher_pusher_proto_rawDescOnce.Do(func() {
		file_pusher_pusher_proto_rawDescData = protoimpl.X.CompressGZIP(file_pusher_pusher_proto_rawDescData)
	})
	return file_pusher_pusher_proto_rawDescData
}

var file_pusher_pusher_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pusher_pusher_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pusher_pusher_proto_goTypes = []interface{}{
	(DelayEnum)(0),                   // 0: com.NimbusIM.proto.pusher.DelayEnum
	(*PusherKafkaRequest)(nil),       // 1: com.NimbusIM.proto.pusher.PusherKafkaRequest
	(*PusherKafkaResponse)(nil),      // 2: com.NimbusIM.proto.pusher.PusherKafkaResponse
	(*PusherKafkaDelayRequest)(nil),  // 3: com.NimbusIM.proto.pusher.PusherKafkaDelayRequest
	(*PusherKafkaDelayResponse)(nil), // 4: com.NimbusIM.proto.pusher.PusherKafkaDelayResponse
}
var file_pusher_pusher_proto_depIdxs = []int32{
	0, // 0: com.NimbusIM.proto.pusher.PusherKafkaDelayRequest.delay:type_name -> com.NimbusIM.proto.pusher.DelayEnum
	1, // 1: com.NimbusIM.proto.pusher.KafkaPusher.PushToKafka:input_type -> com.NimbusIM.proto.pusher.PusherKafkaRequest
	3, // 2: com.NimbusIM.proto.pusher.KafkaPusher.PushToKafkaDelay:input_type -> com.NimbusIM.proto.pusher.PusherKafkaDelayRequest
	2, // 3: com.NimbusIM.proto.pusher.KafkaPusher.PushToKafka:output_type -> com.NimbusIM.proto.pusher.PusherKafkaResponse
	4, // 4: com.NimbusIM.proto.pusher.KafkaPusher.PushToKafkaDelay:output_type -> com.NimbusIM.proto.pusher.PusherKafkaDelayResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pusher_pusher_proto_init() }
func file_pusher_pusher_proto_init() {
	if File_pusher_pusher_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pusher_pusher_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PusherKafkaRequest); i {
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
		file_pusher_pusher_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PusherKafkaResponse); i {
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
		file_pusher_pusher_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PusherKafkaDelayRequest); i {
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
		file_pusher_pusher_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PusherKafkaDelayResponse); i {
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
			RawDescriptor: file_pusher_pusher_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pusher_pusher_proto_goTypes,
		DependencyIndexes: file_pusher_pusher_proto_depIdxs,
		EnumInfos:         file_pusher_pusher_proto_enumTypes,
		MessageInfos:      file_pusher_pusher_proto_msgTypes,
	}.Build()
	File_pusher_pusher_proto = out.File
	file_pusher_pusher_proto_rawDesc = nil
	file_pusher_pusher_proto_goTypes = nil
	file_pusher_pusher_proto_depIdxs = nil
}
