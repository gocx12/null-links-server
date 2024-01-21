// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.6.1
// source: rpc_service/video/video.proto

package video

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

type ParseReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Url   string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *ParseReq) Reset() {
	*x = ParseReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_service_video_video_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ParseReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParseReq) ProtoMessage() {}

func (x *ParseReq) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_service_video_video_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ParseReq.ProtoReflect.Descriptor instead.
func (*ParseReq) Descriptor() ([]byte, []int) {
	return file_rpc_service_video_video_proto_rawDescGZIP(), []int{0}
}

func (x *ParseReq) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *ParseReq) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type ParseResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode int32  `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	StatusMsg  string `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	PlayUrl    int64  `protobuf:"varint,3,opt,name=play_url,json=playUrl,proto3" json:"play_url,omitempty"`
}

func (x *ParseResp) Reset() {
	*x = ParseResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_service_video_video_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ParseResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParseResp) ProtoMessage() {}

func (x *ParseResp) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_service_video_video_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ParseResp.ProtoReflect.Descriptor instead.
func (*ParseResp) Descriptor() ([]byte, []int) {
	return file_rpc_service_video_video_proto_rawDescGZIP(), []int{1}
}

func (x *ParseResp) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *ParseResp) GetStatusMsg() string {
	if x != nil {
		return x.StatusMsg
	}
	return ""
}

func (x *ParseResp) GetPlayUrl() int64 {
	if x != nil {
		return x.PlayUrl
	}
	return 0
}

var File_rpc_service_video_video_proto protoreflect.FileDescriptor

var file_rpc_service_video_video_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x72, 0x70, 0x63, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x69,
	0x64, 0x65, 0x6f, 0x2f, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x22, 0x32, 0x0a, 0x08, 0x50, 0x61, 0x72, 0x73, 0x65, 0x52,
	0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x66, 0x0a, 0x09, 0x50, 0x61,
	0x72, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x4d, 0x73, 0x67, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x5f,
	0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x70, 0x6c, 0x61, 0x79, 0x55,
	0x72, 0x6c, 0x32, 0x3a, 0x0a, 0x0c, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x2a, 0x0a, 0x05, 0x50, 0x61, 0x72, 0x73, 0x65, 0x12, 0x0f, 0x2e, 0x76, 0x69,
	0x64, 0x65, 0x6f, 0x2e, 0x50, 0x61, 0x72, 0x73, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x10, 0x2e, 0x76,
	0x69, 0x64, 0x65, 0x6f, 0x2e, 0x50, 0x61, 0x72, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x42, 0x09,
	0x5a, 0x07, 0x2e, 0x2f, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_rpc_service_video_video_proto_rawDescOnce sync.Once
	file_rpc_service_video_video_proto_rawDescData = file_rpc_service_video_video_proto_rawDesc
)

func file_rpc_service_video_video_proto_rawDescGZIP() []byte {
	file_rpc_service_video_video_proto_rawDescOnce.Do(func() {
		file_rpc_service_video_video_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_service_video_video_proto_rawDescData)
	})
	return file_rpc_service_video_video_proto_rawDescData
}

var file_rpc_service_video_video_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_service_video_video_proto_goTypes = []interface{}{
	(*ParseReq)(nil),  // 0: video.ParseReq
	(*ParseResp)(nil), // 1: video.ParseResp
}
var file_rpc_service_video_video_proto_depIdxs = []int32{
	0, // 0: video.VideoService.Parse:input_type -> video.ParseReq
	1, // 1: video.VideoService.Parse:output_type -> video.ParseResp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rpc_service_video_video_proto_init() }
func file_rpc_service_video_video_proto_init() {
	if File_rpc_service_video_video_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_service_video_video_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ParseReq); i {
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
		file_rpc_service_video_video_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ParseResp); i {
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
			RawDescriptor: file_rpc_service_video_video_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_service_video_video_proto_goTypes,
		DependencyIndexes: file_rpc_service_video_video_proto_depIdxs,
		MessageInfos:      file_rpc_service_video_video_proto_msgTypes,
	}.Build()
	File_rpc_service_video_video_proto = out.File
	file_rpc_service_video_video_proto_rawDesc = nil
	file_rpc_service_video_video_proto_goTypes = nil
	file_rpc_service_video_video_proto_depIdxs = nil
}
