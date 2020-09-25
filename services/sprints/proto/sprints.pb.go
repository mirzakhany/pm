// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: services/sprints/proto/sprints.proto

package sprints

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type SprintStatus int32

const (
	SprintStatus_IN_PROGRESS SprintStatus = 0
	SprintStatus_DONE        SprintStatus = 1
	SprintStatus_ARCHIVED    SprintStatus = 2
)

var SprintStatus_name = map[int32]string{
	0: "IN_PROGRESS",
	1: "DONE",
	2: "ARCHIVED",
}

var SprintStatus_value = map[string]int32{
	"IN_PROGRESS": 0,
	"DONE":        1,
	"ARCHIVED":    2,
}

func (x SprintStatus) String() string {
	return proto.EnumName(SprintStatus_name, int32(x))
}

func (SprintStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_fa7258e795fb1c90, []int{0}
}

type ListSprintsRequest struct {
	PageSize             int32    `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageToken            string   `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListSprintsRequest) Reset()         { *m = ListSprintsRequest{} }
func (m *ListSprintsRequest) String() string { return proto.CompactTextString(m) }
func (*ListSprintsRequest) ProtoMessage()    {}
func (*ListSprintsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fa7258e795fb1c90, []int{0}
}
func (m *ListSprintsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListSprintsRequest.Unmarshal(m, b)
}
func (m *ListSprintsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListSprintsRequest.Marshal(b, m, deterministic)
}
func (m *ListSprintsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListSprintsRequest.Merge(m, src)
}
func (m *ListSprintsRequest) XXX_Size() int {
	return xxx_messageInfo_ListSprintsRequest.Size(m)
}
func (m *ListSprintsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListSprintsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListSprintsRequest proto.InternalMessageInfo

func (m *ListSprintsRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ListSprintsRequest) GetPageToken() string {
	if m != nil {
		return m.PageToken
	}
	return ""
}

type ListSprintsResponse struct {
	Sprints              []*Sprint `protobuf:"bytes,1,rep,name=sprints,proto3" json:"sprints,omitempty"`
	TotalSize            int32     `protobuf:"varint,2,opt,name=total_size,json=totalSize,proto3" json:"total_size,omitempty"`
	NextPageToken        string    `protobuf:"bytes,3,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ListSprintsResponse) Reset()         { *m = ListSprintsResponse{} }
func (m *ListSprintsResponse) String() string { return proto.CompactTextString(m) }
func (*ListSprintsResponse) ProtoMessage()    {}
func (*ListSprintsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_fa7258e795fb1c90, []int{1}
}
func (m *ListSprintsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListSprintsResponse.Unmarshal(m, b)
}
func (m *ListSprintsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListSprintsResponse.Marshal(b, m, deterministic)
}
func (m *ListSprintsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListSprintsResponse.Merge(m, src)
}
func (m *ListSprintsResponse) XXX_Size() int {
	return xxx_messageInfo_ListSprintsResponse.Size(m)
}
func (m *ListSprintsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListSprintsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListSprintsResponse proto.InternalMessageInfo

func (m *ListSprintsResponse) GetSprints() []*Sprint {
	if m != nil {
		return m.Sprints
	}
	return nil
}

func (m *ListSprintsResponse) GetTotalSize() int32 {
	if m != nil {
		return m.TotalSize
	}
	return 0
}

func (m *ListSprintsResponse) GetNextPageToken() string {
	if m != nil {
		return m.NextPageToken
	}
	return ""
}

type GetSprintRequest struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSprintRequest) Reset()         { *m = GetSprintRequest{} }
func (m *GetSprintRequest) String() string { return proto.CompactTextString(m) }
func (*GetSprintRequest) ProtoMessage()    {}
func (*GetSprintRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fa7258e795fb1c90, []int{2}
}
func (m *GetSprintRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSprintRequest.Unmarshal(m, b)
}
func (m *GetSprintRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSprintRequest.Marshal(b, m, deterministic)
}
func (m *GetSprintRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSprintRequest.Merge(m, src)
}
func (m *GetSprintRequest) XXX_Size() int {
	return xxx_messageInfo_GetSprintRequest.Size(m)
}
func (m *GetSprintRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSprintRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSprintRequest proto.InternalMessageInfo

func (m *GetSprintRequest) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type CreateSprintRequest struct {
	Title                string           `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Status               SprintStatus     `protobuf:"varint,2,opt,name=status,proto3,enum=sprints.SprintStatus" json:"status,omitempty"`
	StartAt              *types.Timestamp `protobuf:"bytes,3,opt,name=start_at,json=startAt,proto3" json:"start_at,omitempty"`
	EndAt                *types.Timestamp `protobuf:"bytes,4,opt,name=end_at,json=endAt,proto3" json:"end_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *CreateSprintRequest) Reset()         { *m = CreateSprintRequest{} }
func (m *CreateSprintRequest) String() string { return proto.CompactTextString(m) }
func (*CreateSprintRequest) ProtoMessage()    {}
func (*CreateSprintRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fa7258e795fb1c90, []int{3}
}
func (m *CreateSprintRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateSprintRequest.Unmarshal(m, b)
}
func (m *CreateSprintRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateSprintRequest.Marshal(b, m, deterministic)
}
func (m *CreateSprintRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateSprintRequest.Merge(m, src)
}
func (m *CreateSprintRequest) XXX_Size() int {
	return xxx_messageInfo_CreateSprintRequest.Size(m)
}
func (m *CreateSprintRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateSprintRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateSprintRequest proto.InternalMessageInfo

func (m *CreateSprintRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *CreateSprintRequest) GetStatus() SprintStatus {
	if m != nil {
		return m.Status
	}
	return SprintStatus_IN_PROGRESS
}

func (m *CreateSprintRequest) GetStartAt() *types.Timestamp {
	if m != nil {
		return m.StartAt
	}
	return nil
}

func (m *CreateSprintRequest) GetEndAt() *types.Timestamp {
	if m != nil {
		return m.EndAt
	}
	return nil
}

type UpdateSprintRequest struct {
	Uuid                 string           `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Title                string           `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Status               SprintStatus     `protobuf:"varint,3,opt,name=status,proto3,enum=sprints.SprintStatus" json:"status,omitempty"`
	StartAt              *types.Timestamp `protobuf:"bytes,4,opt,name=start_at,json=startAt,proto3" json:"start_at,omitempty"`
	EndAt                *types.Timestamp `protobuf:"bytes,5,opt,name=end_at,json=endAt,proto3" json:"end_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *UpdateSprintRequest) Reset()         { *m = UpdateSprintRequest{} }
func (m *UpdateSprintRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateSprintRequest) ProtoMessage()    {}
func (*UpdateSprintRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fa7258e795fb1c90, []int{4}
}
func (m *UpdateSprintRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateSprintRequest.Unmarshal(m, b)
}
func (m *UpdateSprintRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateSprintRequest.Marshal(b, m, deterministic)
}
func (m *UpdateSprintRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateSprintRequest.Merge(m, src)
}
func (m *UpdateSprintRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateSprintRequest.Size(m)
}
func (m *UpdateSprintRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateSprintRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateSprintRequest proto.InternalMessageInfo

func (m *UpdateSprintRequest) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *UpdateSprintRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *UpdateSprintRequest) GetStatus() SprintStatus {
	if m != nil {
		return m.Status
	}
	return SprintStatus_IN_PROGRESS
}

func (m *UpdateSprintRequest) GetStartAt() *types.Timestamp {
	if m != nil {
		return m.StartAt
	}
	return nil
}

func (m *UpdateSprintRequest) GetEndAt() *types.Timestamp {
	if m != nil {
		return m.EndAt
	}
	return nil
}

type DeleteSprintRequest struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteSprintRequest) Reset()         { *m = DeleteSprintRequest{} }
func (m *DeleteSprintRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteSprintRequest) ProtoMessage()    {}
func (*DeleteSprintRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fa7258e795fb1c90, []int{5}
}
func (m *DeleteSprintRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteSprintRequest.Unmarshal(m, b)
}
func (m *DeleteSprintRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteSprintRequest.Marshal(b, m, deterministic)
}
func (m *DeleteSprintRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteSprintRequest.Merge(m, src)
}
func (m *DeleteSprintRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteSprintRequest.Size(m)
}
func (m *DeleteSprintRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteSprintRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteSprintRequest proto.InternalMessageInfo

func (m *DeleteSprintRequest) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type Sprint struct {
	Uuid                 string           `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Title                string           `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Status               SprintStatus     `protobuf:"varint,3,opt,name=status,proto3,enum=sprints.SprintStatus" json:"status,omitempty"`
	StartAt              *types.Timestamp `protobuf:"bytes,4,opt,name=start_at,json=startAt,proto3" json:"start_at,omitempty"`
	EndAt                *types.Timestamp `protobuf:"bytes,5,opt,name=end_at,json=endAt,proto3" json:"end_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Sprint) Reset()         { *m = Sprint{} }
func (m *Sprint) String() string { return proto.CompactTextString(m) }
func (*Sprint) ProtoMessage()    {}
func (*Sprint) Descriptor() ([]byte, []int) {
	return fileDescriptor_fa7258e795fb1c90, []int{6}
}
func (m *Sprint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Sprint.Unmarshal(m, b)
}
func (m *Sprint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Sprint.Marshal(b, m, deterministic)
}
func (m *Sprint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Sprint.Merge(m, src)
}
func (m *Sprint) XXX_Size() int {
	return xxx_messageInfo_Sprint.Size(m)
}
func (m *Sprint) XXX_DiscardUnknown() {
	xxx_messageInfo_Sprint.DiscardUnknown(m)
}

var xxx_messageInfo_Sprint proto.InternalMessageInfo

func (m *Sprint) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Sprint) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Sprint) GetStatus() SprintStatus {
	if m != nil {
		return m.Status
	}
	return SprintStatus_IN_PROGRESS
}

func (m *Sprint) GetStartAt() *types.Timestamp {
	if m != nil {
		return m.StartAt
	}
	return nil
}

func (m *Sprint) GetEndAt() *types.Timestamp {
	if m != nil {
		return m.EndAt
	}
	return nil
}

func init() {
	proto.RegisterEnum("sprints.SprintStatus", SprintStatus_name, SprintStatus_value)
	proto.RegisterType((*ListSprintsRequest)(nil), "sprints.ListSprintsRequest")
	proto.RegisterType((*ListSprintsResponse)(nil), "sprints.ListSprintsResponse")
	proto.RegisterType((*GetSprintRequest)(nil), "sprints.GetSprintRequest")
	proto.RegisterType((*CreateSprintRequest)(nil), "sprints.CreateSprintRequest")
	proto.RegisterType((*UpdateSprintRequest)(nil), "sprints.UpdateSprintRequest")
	proto.RegisterType((*DeleteSprintRequest)(nil), "sprints.DeleteSprintRequest")
	proto.RegisterType((*Sprint)(nil), "sprints.Sprint")
}

func init() {
	proto.RegisterFile("services/sprints/proto/sprints.proto", fileDescriptor_fa7258e795fb1c90)
}

var fileDescriptor_fa7258e795fb1c90 = []byte{
	// 606 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xdc, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xf9, 0x6b, 0x3c, 0x49, 0x89, 0x99, 0x94, 0x36, 0x38, 0xad, 0x88, 0x2c, 0x54, 0xa5,
	0x91, 0x88, 0xd5, 0x20, 0x84, 0xc4, 0x2d, 0x34, 0x51, 0xa9, 0x84, 0xda, 0xc8, 0x69, 0x39, 0x54,
	0x48, 0xc1, 0x25, 0x4b, 0x64, 0x91, 0xd8, 0x26, 0x3b, 0xa9, 0xa0, 0x88, 0x0b, 0x17, 0x1e, 0x80,
	0x57, 0xe2, 0xc6, 0x11, 0x0e, 0x3c, 0x00, 0x0f, 0x82, 0xbc, 0xf6, 0xa6, 0x6e, 0x62, 0x89, 0xc2,
	0x91, 0x9b, 0x67, 0x67, 0xe6, 0x9b, 0xef, 0xd3, 0xf8, 0x1b, 0xb8, 0xc7, 0xd9, 0xf4, 0xdc, 0x79,
	0xc5, 0xb8, 0xc9, 0xfd, 0xa9, 0xe3, 0x12, 0x37, 0xfd, 0xa9, 0x47, 0x9e, 0x8c, 0x9a, 0x22, 0xc2,
	0x95, 0x28, 0xd4, 0x37, 0x47, 0x9e, 0x37, 0x1a, 0x33, 0xd3, 0xf6, 0x1d, 0xd3, 0x76, 0x5d, 0x8f,
	0x6c, 0x72, 0x3c, 0x37, 0x2a, 0xd3, 0xab, 0x51, 0x56, 0x44, 0x67, 0xb3, 0xd7, 0x26, 0x9b, 0xf8,
	0xf4, 0x3e, 0x4a, 0xde, 0x5d, 0x4c, 0x92, 0x33, 0x61, 0x9c, 0xec, 0x89, 0x1f, 0x16, 0x18, 0x3d,
	0xc0, 0x67, 0x0e, 0xa7, 0x7e, 0x38, 0xca, 0x62, 0x6f, 0x67, 0x8c, 0x13, 0x56, 0x41, 0xf5, 0xed,
	0x11, 0x1b, 0x70, 0xe7, 0x82, 0x55, 0x94, 0x9a, 0x52, 0xcf, 0x5a, 0xf9, 0xe0, 0xa1, 0xef, 0x5c,
	0x30, 0xdc, 0x02, 0x10, 0x49, 0xf2, 0xde, 0x30, 0xb7, 0x92, 0xaa, 0x29, 0x75, 0xd5, 0x12, 0xe5,
	0xc7, 0xc1, 0x83, 0xf1, 0x59, 0x81, 0xf2, 0x15, 0x48, 0xee, 0x7b, 0x2e, 0x67, 0xb8, 0x03, 0x52,
	0x50, 0x45, 0xa9, 0xa5, 0xeb, 0x85, 0x56, 0xa9, 0x29, 0xf5, 0x86, 0xa5, 0x96, 0xcc, 0x07, 0x13,
	0xc8, 0x23, 0x7b, 0x1c, 0xce, 0x4f, 0x89, 0xf9, 0xaa, 0x78, 0x11, 0x04, 0xb6, 0xa1, 0xe4, 0xb2,
	0x77, 0x34, 0x88, 0xb1, 0x48, 0x0b, 0x16, 0xab, 0xc1, 0x73, 0x6f, 0xce, 0x64, 0x1b, 0xb4, 0x7d,
	0x16, 0xf1, 0x90, 0xca, 0x10, 0x32, 0xb3, 0x99, 0x33, 0x14, 0xa2, 0x54, 0x4b, 0x7c, 0x1b, 0x5f,
	0x15, 0x28, 0xef, 0x4d, 0x99, 0x4d, 0xec, 0x6a, 0xed, 0x1a, 0x64, 0xc9, 0xa1, 0x31, 0x8b, 0x8a,
	0xc3, 0x00, 0xef, 0x43, 0x8e, 0x93, 0x4d, 0x33, 0x2e, 0x88, 0xdd, 0x6c, 0xdd, 0x5e, 0x90, 0xd1,
	0x17, 0x49, 0x2b, 0x2a, 0xc2, 0x87, 0x90, 0xe7, 0x64, 0x4f, 0x69, 0x60, 0x93, 0x60, 0x59, 0x68,
	0xe9, 0xcd, 0x70, 0x29, 0x4d, 0xb9, 0x94, 0xe6, 0xb1, 0x5c, 0x8a, 0xb5, 0x22, 0x6a, 0xdb, 0x84,
	0xbb, 0x90, 0x63, 0xee, 0x30, 0x68, 0xca, 0xfc, 0xb1, 0x29, 0xcb, 0xdc, 0x61, 0x9b, 0x8c, 0x9f,
	0x0a, 0x94, 0x4f, 0xfc, 0xe1, 0x92, 0x8c, 0x04, 0xc9, 0x97, 0xd2, 0x52, 0xc9, 0xd2, 0xd2, 0x7f,
	0x2b, 0x2d, 0xf3, 0x2f, 0xd2, 0xb2, 0xd7, 0x95, 0xb6, 0x03, 0xe5, 0x0e, 0x1b, 0xb3, 0x6b, 0x28,
	0x33, 0xbe, 0x29, 0x90, 0x0b, 0xab, 0xfe, 0x03, 0xe1, 0x8d, 0x47, 0x50, 0x8c, 0x33, 0xc0, 0x12,
	0x14, 0x0e, 0x0e, 0x07, 0x3d, 0xeb, 0x68, 0xdf, 0xea, 0xf6, 0xfb, 0xda, 0x0d, 0xcc, 0x43, 0xa6,
	0x73, 0x74, 0xd8, 0xd5, 0x14, 0x2c, 0x42, 0xbe, 0x6d, 0xed, 0x3d, 0x3d, 0x78, 0xde, 0xed, 0x68,
	0xa9, 0xd6, 0x8f, 0x34, 0xac, 0x46, 0x9d, 0xe1, 0xad, 0xc1, 0x53, 0x28, 0xc4, 0x6c, 0x89, 0xd5,
	0xb9, 0xc4, 0x65, 0xff, 0xeb, 0x9b, 0xc9, 0xc9, 0xd0, 0xc9, 0x86, 0xf6, 0xe9, 0xfb, 0xaf, 0x2f,
	0x29, 0xc0, 0xbc, 0x3c, 0x58, 0xd8, 0x03, 0x75, 0xee, 0x34, 0xbc, 0x33, 0x6f, 0x5e, 0x74, 0x9f,
	0xbe, 0x68, 0x79, 0x63, 0x43, 0x40, 0xdd, 0xc2, 0xd2, 0xfc, 0x12, 0x7e, 0x08, 0xd6, 0xf4, 0x11,
	0x4f, 0xa0, 0x18, 0xb7, 0x24, 0x5e, 0x32, 0x4a, 0x70, 0xea, 0x32, 0xee, 0xba, 0xc0, 0xd5, 0x8c,
	0x82, 0x79, 0xbe, 0x2b, 0xa1, 0x1f, 0x2b, 0x0d, 0x7c, 0x01, 0xc5, 0xb8, 0x45, 0x62, 0xb0, 0x09,
	0xce, 0x59, 0x86, 0xdd, 0x12, 0xb0, 0x1b, 0x3a, 0xc6, 0x60, 0x23, 0xc6, 0x01, 0xfa, 0x4b, 0x28,
	0xc6, 0x7f, 0xd3, 0x18, 0x7a, 0xc2, 0xdf, 0xab, 0xaf, 0x2f, 0xad, 0xbf, 0x1b, 0x5c, 0x6e, 0x43,
	0x17, 0x43, 0xd6, 0x1a, 0x09, 0x43, 0x9e, 0xa8, 0xa7, 0xf2, 0x48, 0x9e, 0xe5, 0x44, 0xdb, 0x83,
	0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x2b, 0xe7, 0xd2, 0xe2, 0x4d, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SprintServiceClient is the client API for SprintService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SprintServiceClient interface {
	// List Sprints
	ListSprints(ctx context.Context, in *ListSprintsRequest, opts ...grpc.CallOption) (*ListSprintsResponse, error)
	// Get Sprint
	GetSprint(ctx context.Context, in *GetSprintRequest, opts ...grpc.CallOption) (*Sprint, error)
	// Create Sprint object request
	CreateSprint(ctx context.Context, in *CreateSprintRequest, opts ...grpc.CallOption) (*Sprint, error)
	// Update Sprint object request
	UpdateSprint(ctx context.Context, in *UpdateSprintRequest, opts ...grpc.CallOption) (*Sprint, error)
	// Delete Sprint object request
	DeleteSprint(ctx context.Context, in *DeleteSprintRequest, opts ...grpc.CallOption) (*types.Empty, error)
}

type sprintServiceClient struct {
	cc *grpc.ClientConn
}

func NewSprintServiceClient(cc *grpc.ClientConn) SprintServiceClient {
	return &sprintServiceClient{cc}
}

func (c *sprintServiceClient) ListSprints(ctx context.Context, in *ListSprintsRequest, opts ...grpc.CallOption) (*ListSprintsResponse, error) {
	out := new(ListSprintsResponse)
	err := c.cc.Invoke(ctx, "/sprints.SprintService/ListSprints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sprintServiceClient) GetSprint(ctx context.Context, in *GetSprintRequest, opts ...grpc.CallOption) (*Sprint, error) {
	out := new(Sprint)
	err := c.cc.Invoke(ctx, "/sprints.SprintService/GetSprint", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sprintServiceClient) CreateSprint(ctx context.Context, in *CreateSprintRequest, opts ...grpc.CallOption) (*Sprint, error) {
	out := new(Sprint)
	err := c.cc.Invoke(ctx, "/sprints.SprintService/CreateSprint", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sprintServiceClient) UpdateSprint(ctx context.Context, in *UpdateSprintRequest, opts ...grpc.CallOption) (*Sprint, error) {
	out := new(Sprint)
	err := c.cc.Invoke(ctx, "/sprints.SprintService/UpdateSprint", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sprintServiceClient) DeleteSprint(ctx context.Context, in *DeleteSprintRequest, opts ...grpc.CallOption) (*types.Empty, error) {
	out := new(types.Empty)
	err := c.cc.Invoke(ctx, "/sprints.SprintService/DeleteSprint", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SprintServiceServer is the server API for SprintService service.
type SprintServiceServer interface {
	// List Sprints
	ListSprints(context.Context, *ListSprintsRequest) (*ListSprintsResponse, error)
	// Get Sprint
	GetSprint(context.Context, *GetSprintRequest) (*Sprint, error)
	// Create Sprint object request
	CreateSprint(context.Context, *CreateSprintRequest) (*Sprint, error)
	// Update Sprint object request
	UpdateSprint(context.Context, *UpdateSprintRequest) (*Sprint, error)
	// Delete Sprint object request
	DeleteSprint(context.Context, *DeleteSprintRequest) (*types.Empty, error)
}

// UnimplementedSprintServiceServer can be embedded to have forward compatible implementations.
type UnimplementedSprintServiceServer struct {
}

func (*UnimplementedSprintServiceServer) ListSprints(ctx context.Context, req *ListSprintsRequest) (*ListSprintsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSprints not implemented")
}
func (*UnimplementedSprintServiceServer) GetSprint(ctx context.Context, req *GetSprintRequest) (*Sprint, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSprint not implemented")
}
func (*UnimplementedSprintServiceServer) CreateSprint(ctx context.Context, req *CreateSprintRequest) (*Sprint, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSprint not implemented")
}
func (*UnimplementedSprintServiceServer) UpdateSprint(ctx context.Context, req *UpdateSprintRequest) (*Sprint, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSprint not implemented")
}
func (*UnimplementedSprintServiceServer) DeleteSprint(ctx context.Context, req *DeleteSprintRequest) (*types.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSprint not implemented")
}

func RegisterSprintServiceServer(s *grpc.Server, srv SprintServiceServer) {
	s.RegisterService(&_SprintService_serviceDesc, srv)
}

func _SprintService_ListSprints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSprintsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SprintServiceServer).ListSprints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sprints.SprintService/ListSprints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SprintServiceServer).ListSprints(ctx, req.(*ListSprintsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SprintService_GetSprint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSprintRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SprintServiceServer).GetSprint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sprints.SprintService/GetSprint",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SprintServiceServer).GetSprint(ctx, req.(*GetSprintRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SprintService_CreateSprint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSprintRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SprintServiceServer).CreateSprint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sprints.SprintService/CreateSprint",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SprintServiceServer).CreateSprint(ctx, req.(*CreateSprintRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SprintService_UpdateSprint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateSprintRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SprintServiceServer).UpdateSprint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sprints.SprintService/UpdateSprint",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SprintServiceServer).UpdateSprint(ctx, req.(*UpdateSprintRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SprintService_DeleteSprint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSprintRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SprintServiceServer).DeleteSprint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sprints.SprintService/DeleteSprint",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SprintServiceServer).DeleteSprint(ctx, req.(*DeleteSprintRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SprintService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sprints.SprintService",
	HandlerType: (*SprintServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListSprints",
			Handler:    _SprintService_ListSprints_Handler,
		},
		{
			MethodName: "GetSprint",
			Handler:    _SprintService_GetSprint_Handler,
		},
		{
			MethodName: "CreateSprint",
			Handler:    _SprintService_CreateSprint_Handler,
		},
		{
			MethodName: "UpdateSprint",
			Handler:    _SprintService_UpdateSprint_Handler,
		},
		{
			MethodName: "DeleteSprint",
			Handler:    _SprintService_DeleteSprint_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/sprints/proto/sprints.proto",
}
