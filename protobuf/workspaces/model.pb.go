// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.8.0
// source: protobuf/workspaces/model.proto

package workspaces

import (
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Workspace struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Uuid      string               `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Title     string               `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Domain    string               `protobuf:"bytes,4,opt,name=domain,proto3" json:"domain,omitempty"`
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamp.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Workspace) Reset() {
	*x = Workspace{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protobuf_workspaces_model_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Workspace) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Workspace) ProtoMessage() {}

func (x *Workspace) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_workspaces_model_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Workspace.ProtoReflect.Descriptor instead.
func (*Workspace) Descriptor() ([]byte, []int) {
	return file_protobuf_workspaces_model_proto_rawDescGZIP(), []int{0}
}

func (x *Workspace) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Workspace) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *Workspace) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Workspace) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *Workspace) GetCreatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Workspace) GetUpdatedAt() *timestamp.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

var File_protobuf_workspaces_model_proto protoreflect.FileDescriptor

var file_protobuf_workspaces_model_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x73, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x56, 0x31, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xd3, 0x01, 0x0a, 0x09, 0x57, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75,
	0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x42, 0x20, 0x5a, 0x1e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x3b, 0x77, 0x6f,
	0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protobuf_workspaces_model_proto_rawDescOnce sync.Once
	file_protobuf_workspaces_model_proto_rawDescData = file_protobuf_workspaces_model_proto_rawDesc
)

func file_protobuf_workspaces_model_proto_rawDescGZIP() []byte {
	file_protobuf_workspaces_model_proto_rawDescOnce.Do(func() {
		file_protobuf_workspaces_model_proto_rawDescData = protoimpl.X.CompressGZIP(file_protobuf_workspaces_model_proto_rawDescData)
	})
	return file_protobuf_workspaces_model_proto_rawDescData
}

var file_protobuf_workspaces_model_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_protobuf_workspaces_model_proto_goTypes = []interface{}{
	(*Workspace)(nil),           // 0: workspacesV1.Workspace
	(*timestamp.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_protobuf_workspaces_model_proto_depIdxs = []int32{
	1, // 0: workspacesV1.Workspace.created_at:type_name -> google.protobuf.Timestamp
	1, // 1: workspacesV1.Workspace.updated_at:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_protobuf_workspaces_model_proto_init() }
func file_protobuf_workspaces_model_proto_init() {
	if File_protobuf_workspaces_model_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protobuf_workspaces_model_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Workspace); i {
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
			RawDescriptor: file_protobuf_workspaces_model_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protobuf_workspaces_model_proto_goTypes,
		DependencyIndexes: file_protobuf_workspaces_model_proto_depIdxs,
		MessageInfos:      file_protobuf_workspaces_model_proto_msgTypes,
	}.Build()
	File_protobuf_workspaces_model_proto = out.File
	file_protobuf_workspaces_model_proto_rawDesc = nil
	file_protobuf_workspaces_model_proto_goTypes = nil
	file_protobuf_workspaces_model_proto_depIdxs = nil
}
