// proto version

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.14.0
// source: module/book/bookPb/bookPb.proto

package bookPb

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

type FindBooksInIdsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids []uint64 `protobuf:"varint,1,rep,packed,name=ids,proto3" json:"ids,omitempty"`
}

func (x *FindBooksInIdsReq) Reset() {
	*x = FindBooksInIdsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_module_book_bookPb_bookPb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindBooksInIdsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindBooksInIdsReq) ProtoMessage() {}

func (x *FindBooksInIdsReq) ProtoReflect() protoreflect.Message {
	mi := &file_module_book_bookPb_bookPb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindBooksInIdsReq.ProtoReflect.Descriptor instead.
func (*FindBooksInIdsReq) Descriptor() ([]byte, []int) {
	return file_module_book_bookPb_bookPb_proto_rawDescGZIP(), []int{0}
}

func (x *FindBooksInIdsReq) GetIds() []uint64 {
	if x != nil {
		return x.Ids
	}
	return nil
}

type FindBooksInIdsRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Book []*Book `protobuf:"bytes,1,rep,name=book,proto3" json:"book,omitempty"`
}

func (x *FindBooksInIdsRes) Reset() {
	*x = FindBooksInIdsRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_module_book_bookPb_bookPb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindBooksInIdsRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindBooksInIdsRes) ProtoMessage() {}

func (x *FindBooksInIdsRes) ProtoReflect() protoreflect.Message {
	mi := &file_module_book_bookPb_bookPb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindBooksInIdsRes.ProtoReflect.Descriptor instead.
func (*FindBooksInIdsRes) Descriptor() ([]byte, []int) {
	return file_module_book_bookPb_bookPb_proto_rawDescGZIP(), []int{1}
}

func (x *FindBooksInIdsRes) GetBook() []*Book {
	if x != nil {
		return x.Book
	}
	return nil
}

type Book struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             uint64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title          string  `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Price          uint64  `protobuf:"varint,3,opt,name=price,proto3" json:"price,omitempty"`
	FilePath       string  `protobuf:"bytes,4,opt,name=file_path,json=filePath,proto3" json:"file_path,omitempty"`
	CoverImagePath string  `protobuf:"bytes,5,opt,name=cover_image_path,json=coverImagePath,proto3" json:"cover_image_path,omitempty"`
	AuthorName     string  `protobuf:"bytes,6,opt,name=author_name,json=authorName,proto3" json:"author_name,omitempty"`
	Tags           []*Tags `protobuf:"bytes,7,rep,name=tags,proto3" json:"tags,omitempty"`
}

func (x *Book) Reset() {
	*x = Book{}
	if protoimpl.UnsafeEnabled {
		mi := &file_module_book_bookPb_bookPb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Book) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Book) ProtoMessage() {}

func (x *Book) ProtoReflect() protoreflect.Message {
	mi := &file_module_book_bookPb_bookPb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Book.ProtoReflect.Descriptor instead.
func (*Book) Descriptor() ([]byte, []int) {
	return file_module_book_bookPb_bookPb_proto_rawDescGZIP(), []int{2}
}

func (x *Book) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Book) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Book) GetPrice() uint64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Book) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

func (x *Book) GetCoverImagePath() string {
	if x != nil {
		return x.CoverImagePath
	}
	return ""
}

func (x *Book) GetAuthorName() string {
	if x != nil {
		return x.AuthorName
	}
	return ""
}

func (x *Book) GetTags() []*Tags {
	if x != nil {
		return x.Tags
	}
	return nil
}

type Tags struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Tags) Reset() {
	*x = Tags{}
	if protoimpl.UnsafeEnabled {
		mi := &file_module_book_bookPb_bookPb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tags) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tags) ProtoMessage() {}

func (x *Tags) ProtoReflect() protoreflect.Message {
	mi := &file_module_book_bookPb_bookPb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tags.ProtoReflect.Descriptor instead.
func (*Tags) Descriptor() ([]byte, []int) {
	return file_module_book_bookPb_bookPb_proto_rawDescGZIP(), []int{3}
}

func (x *Tags) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Tags) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_module_book_bookPb_bookPb_proto protoreflect.FileDescriptor

var file_module_book_bookPb_bookPb_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x2f, 0x62, 0x6f, 0x6f, 0x6b, 0x2f, 0x62, 0x6f,
	0x6f, 0x6b, 0x50, 0x62, 0x2f, 0x62, 0x6f, 0x6f, 0x6b, 0x50, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x06, 0x62, 0x6f, 0x6f, 0x6b, 0x50, 0x62, 0x22, 0x25, 0x0a, 0x11, 0x46, 0x69, 0x6e,
	0x64, 0x42, 0x6f, 0x6f, 0x6b, 0x73, 0x49, 0x6e, 0x49, 0x64, 0x73, 0x52, 0x65, 0x71, 0x12, 0x10,
	0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x04, 0x52, 0x03, 0x69, 0x64, 0x73,
	0x22, 0x35, 0x0a, 0x11, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x6f, 0x6f, 0x6b, 0x73, 0x49, 0x6e, 0x49,
	0x64, 0x73, 0x52, 0x65, 0x73, 0x12, 0x20, 0x0a, 0x04, 0x62, 0x6f, 0x6f, 0x6b, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x62, 0x6f, 0x6f, 0x6b, 0x50, 0x62, 0x2e, 0x42, 0x6f, 0x6f,
	0x6b, 0x52, 0x04, 0x62, 0x6f, 0x6f, 0x6b, 0x22, 0xcc, 0x01, 0x0a, 0x04, 0x42, 0x6f, 0x6f, 0x6b,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x66, 0x69, 0x6c, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x28, 0x0a, 0x10, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x50,
	0x61, 0x74, 0x68, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x07, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x62, 0x6f, 0x6f, 0x6b, 0x50, 0x62, 0x2e, 0x54, 0x61, 0x67, 0x73,
	0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x22, 0x2a, 0x0a, 0x04, 0x54, 0x61, 0x67, 0x73, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x32, 0x59, 0x0a, 0x0f, 0x42, 0x6f, 0x6f, 0x6b, 0x47, 0x72, 0x70, 0x63, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x46, 0x0a, 0x0e, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x6f, 0x6f,
	0x6b, 0x73, 0x49, 0x6e, 0x49, 0x64, 0x73, 0x12, 0x19, 0x2e, 0x62, 0x6f, 0x6f, 0x6b, 0x50, 0x62,
	0x2e, 0x46, 0x69, 0x6e, 0x64, 0x42, 0x6f, 0x6f, 0x6b, 0x73, 0x49, 0x6e, 0x49, 0x64, 0x73, 0x52,
	0x65, 0x71, 0x1a, 0x19, 0x2e, 0x62, 0x6f, 0x6f, 0x6b, 0x50, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x64,
	0x42, 0x6f, 0x6f, 0x6b, 0x73, 0x49, 0x6e, 0x49, 0x64, 0x73, 0x52, 0x65, 0x73, 0x42, 0x30, 0x5a,
	0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x72, 0x69, 0x74,
	0x41, 0x73, 0x61, 0x77, 0x61, 0x6e, 0x69, 0x72, 0x61, 0x6d, 0x6f, 0x6c, 0x2f, 0x6d, 0x6f, 0x64,
	0x75, 0x6c, 0x65, 0x2f, 0x62, 0x6f, 0x6f, 0x6b, 0x2f, 0x62, 0x6f, 0x6f, 0x6b, 0x50, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_module_book_bookPb_bookPb_proto_rawDescOnce sync.Once
	file_module_book_bookPb_bookPb_proto_rawDescData = file_module_book_bookPb_bookPb_proto_rawDesc
)

func file_module_book_bookPb_bookPb_proto_rawDescGZIP() []byte {
	file_module_book_bookPb_bookPb_proto_rawDescOnce.Do(func() {
		file_module_book_bookPb_bookPb_proto_rawDescData = protoimpl.X.CompressGZIP(file_module_book_bookPb_bookPb_proto_rawDescData)
	})
	return file_module_book_bookPb_bookPb_proto_rawDescData
}

var file_module_book_bookPb_bookPb_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_module_book_bookPb_bookPb_proto_goTypes = []any{
	(*FindBooksInIdsReq)(nil), // 0: bookPb.FindBooksInIdsReq
	(*FindBooksInIdsRes)(nil), // 1: bookPb.FindBooksInIdsRes
	(*Book)(nil),              // 2: bookPb.Book
	(*Tags)(nil),              // 3: bookPb.Tags
}
var file_module_book_bookPb_bookPb_proto_depIdxs = []int32{
	2, // 0: bookPb.FindBooksInIdsRes.book:type_name -> bookPb.Book
	3, // 1: bookPb.Book.tags:type_name -> bookPb.Tags
	0, // 2: bookPb.BookGrpcService.FindBooksInIds:input_type -> bookPb.FindBooksInIdsReq
	1, // 3: bookPb.BookGrpcService.FindBooksInIds:output_type -> bookPb.FindBooksInIdsRes
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_module_book_bookPb_bookPb_proto_init() }
func file_module_book_bookPb_bookPb_proto_init() {
	if File_module_book_bookPb_bookPb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_module_book_bookPb_bookPb_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*FindBooksInIdsReq); i {
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
		file_module_book_bookPb_bookPb_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*FindBooksInIdsRes); i {
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
		file_module_book_bookPb_bookPb_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Book); i {
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
		file_module_book_bookPb_bookPb_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Tags); i {
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
			RawDescriptor: file_module_book_bookPb_bookPb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_module_book_bookPb_bookPb_proto_goTypes,
		DependencyIndexes: file_module_book_bookPb_bookPb_proto_depIdxs,
		MessageInfos:      file_module_book_bookPb_bookPb_proto_msgTypes,
	}.Build()
	File_module_book_bookPb_bookPb_proto = out.File
	file_module_book_bookPb_bookPb_proto_rawDesc = nil
	file_module_book_bookPb_bookPb_proto_goTypes = nil
	file_module_book_bookPb_bookPb_proto_depIdxs = nil
}