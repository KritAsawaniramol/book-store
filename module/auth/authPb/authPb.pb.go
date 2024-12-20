// proto version

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.1
// source: module/auth/authPb/authPb.proto

// Package name

package authPb

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

type AccessTokenSearchReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken string `protobuf:"bytes,1,opt,name=accessToken,proto3" json:"accessToken,omitempty"`
}

func (x *AccessTokenSearchReq) Reset() {
	*x = AccessTokenSearchReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_module_auth_authPb_authPb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessTokenSearchReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessTokenSearchReq) ProtoMessage() {}

func (x *AccessTokenSearchReq) ProtoReflect() protoreflect.Message {
	mi := &file_module_auth_authPb_authPb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessTokenSearchReq.ProtoReflect.Descriptor instead.
func (*AccessTokenSearchReq) Descriptor() ([]byte, []int) {
	return file_module_auth_authPb_authPb_proto_rawDescGZIP(), []int{0}
}

func (x *AccessTokenSearchReq) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

type AccessTokenSearchRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsValid bool `protobuf:"varint,1,opt,name=isValid,proto3" json:"isValid,omitempty"`
}

func (x *AccessTokenSearchRes) Reset() {
	*x = AccessTokenSearchRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_module_auth_authPb_authPb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessTokenSearchRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessTokenSearchRes) ProtoMessage() {}

func (x *AccessTokenSearchRes) ProtoReflect() protoreflect.Message {
	mi := &file_module_auth_authPb_authPb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessTokenSearchRes.ProtoReflect.Descriptor instead.
func (*AccessTokenSearchRes) Descriptor() ([]byte, []int) {
	return file_module_auth_authPb_authPb_proto_rawDescGZIP(), []int{1}
}

func (x *AccessTokenSearchRes) GetIsValid() bool {
	if x != nil {
		return x.IsValid
	}
	return false
}

var File_module_auth_authPb_authPb_proto protoreflect.FileDescriptor

var file_module_auth_authPb_authPb_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x61, 0x75,
	0x74, 0x68, 0x50, 0x62, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x50, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x06, 0x75, 0x73, 0x65, 0x72, 0x50, 0x62, 0x22, 0x38, 0x0a, 0x14, 0x41, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65,
	0x71, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x22, 0x30, 0x0a, 0x14, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x69,
	0x73, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73,
	0x56, 0x61, 0x6c, 0x69, 0x64, 0x32, 0x62, 0x0a, 0x0f, 0x41, 0x75, 0x74, 0x68, 0x47, 0x72, 0x70,
	0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4f, 0x0a, 0x11, 0x41, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x1c, 0x2e,
	0x75, 0x73, 0x65, 0x72, 0x50, 0x62, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x71, 0x1a, 0x1c, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x50, 0x62, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x72, 0x69, 0x74, 0x41, 0x73, 0x61, 0x77,
	0x61, 0x6e, 0x69, 0x72, 0x61, 0x6d, 0x6f, 0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2f,
	0x61, 0x75, 0x74, 0x68, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x50, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_module_auth_authPb_authPb_proto_rawDescOnce sync.Once
	file_module_auth_authPb_authPb_proto_rawDescData = file_module_auth_authPb_authPb_proto_rawDesc
)

func file_module_auth_authPb_authPb_proto_rawDescGZIP() []byte {
	file_module_auth_authPb_authPb_proto_rawDescOnce.Do(func() {
		file_module_auth_authPb_authPb_proto_rawDescData = protoimpl.X.CompressGZIP(file_module_auth_authPb_authPb_proto_rawDescData)
	})
	return file_module_auth_authPb_authPb_proto_rawDescData
}

var file_module_auth_authPb_authPb_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_module_auth_authPb_authPb_proto_goTypes = []any{
	(*AccessTokenSearchReq)(nil), // 0: userPb.AccessTokenSearchReq
	(*AccessTokenSearchRes)(nil), // 1: userPb.AccessTokenSearchRes
}
var file_module_auth_authPb_authPb_proto_depIdxs = []int32{
	0, // 0: userPb.AuthGrpcService.AccessTokenSearch:input_type -> userPb.AccessTokenSearchReq
	1, // 1: userPb.AuthGrpcService.AccessTokenSearch:output_type -> userPb.AccessTokenSearchRes
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_module_auth_authPb_authPb_proto_init() }
func file_module_auth_authPb_authPb_proto_init() {
	if File_module_auth_authPb_authPb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_module_auth_authPb_authPb_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*AccessTokenSearchReq); i {
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
		file_module_auth_authPb_authPb_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*AccessTokenSearchRes); i {
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
			RawDescriptor: file_module_auth_authPb_authPb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_module_auth_authPb_authPb_proto_goTypes,
		DependencyIndexes: file_module_auth_authPb_authPb_proto_depIdxs,
		MessageInfos:      file_module_auth_authPb_authPb_proto_msgTypes,
	}.Build()
	File_module_auth_authPb_authPb_proto = out.File
	file_module_auth_authPb_authPb_proto_rawDesc = nil
	file_module_auth_authPb_authPb_proto_goTypes = nil
	file_module_auth_authPb_authPb_proto_depIdxs = nil
}
