// Code generated by protoc-gen-go. DO NOT EDIT.
// source: tg_account/tg_account.proto

package tg_account

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type TgAccount struct {
	ID                   uint64               `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	TgID                 string               `protobuf:"bytes,2,opt,name=TgID,proto3" json:"TgID,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,3,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *TgAccount) Reset()         { *m = TgAccount{} }
func (m *TgAccount) String() string { return proto.CompactTextString(m) }
func (*TgAccount) ProtoMessage()    {}
func (*TgAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_8b2ac9b320d8f54d, []int{0}
}

func (m *TgAccount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TgAccount.Unmarshal(m, b)
}
func (m *TgAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TgAccount.Marshal(b, m, deterministic)
}
func (m *TgAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TgAccount.Merge(m, src)
}
func (m *TgAccount) XXX_Size() int {
	return xxx_messageInfo_TgAccount.Size(m)
}
func (m *TgAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_TgAccount.DiscardUnknown(m)
}

var xxx_messageInfo_TgAccount proto.InternalMessageInfo

func (m *TgAccount) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *TgAccount) GetTgID() string {
	if m != nil {
		return m.TgID
	}
	return ""
}

func (m *TgAccount) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func init() {
	proto.RegisterType((*TgAccount)(nil), "tg_account.TgAccount")
}

func init() {
	proto.RegisterFile("tg_account/tg_account.proto", fileDescriptor_8b2ac9b320d8f54d)
}

var fileDescriptor_8b2ac9b320d8f54d = []byte{
	// 193 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0x31, 0x0f, 0x82, 0x30,
	0x10, 0x85, 0x53, 0x24, 0x26, 0xd4, 0xc4, 0xa1, 0x13, 0xc1, 0x41, 0xe2, 0xc4, 0xd4, 0x26, 0xe8,
	0xe0, 0x8a, 0xb2, 0xb0, 0x12, 0x26, 0x17, 0x53, 0xb0, 0xd6, 0x26, 0x40, 0x9b, 0x7a, 0x1d, 0xf4,
	0xd7, 0x9b, 0x94, 0x10, 0xdc, 0xde, 0xbd, 0xfb, 0x72, 0x5f, 0x0e, 0xef, 0x40, 0xde, 0x79, 0xd7,
	0x69, 0x37, 0x02, 0x5b, 0x22, 0x35, 0x56, 0x83, 0x26, 0x78, 0x69, 0x92, 0xbd, 0xd4, 0x5a, 0xf6,
	0x82, 0xf9, 0x4d, 0xeb, 0x9e, 0x0c, 0xd4, 0x20, 0xde, 0xc0, 0x07, 0x33, 0xc1, 0x07, 0x85, 0xa3,
	0x46, 0x16, 0x13, 0x4d, 0xb6, 0x38, 0xa8, 0xca, 0x18, 0xa5, 0x28, 0x0b, 0xeb, 0xa0, 0x2a, 0x09,
	0xc1, 0x61, 0x23, 0xab, 0x32, 0x0e, 0x52, 0x94, 0x45, 0xb5, 0xcf, 0xe4, 0x8c, 0xa3, 0xab, 0x15,
	0x1c, 0xc4, 0xa3, 0x80, 0x78, 0x95, 0xa2, 0x6c, 0x93, 0x27, 0x74, 0xb2, 0xd0, 0xd9, 0x42, 0x9b,
	0xd9, 0x52, 0x2f, 0xf0, 0xe5, 0x74, 0xcb, 0xa5, 0x82, 0x97, 0x6b, 0x69, 0xa7, 0x07, 0xf6, 0xe1,
	0xd6, 0x7d, 0x19, 0x37, 0x86, 0xa9, 0x11, 0x84, 0x1d, 0x79, 0xef, 0x07, 0x7f, 0xe2, 0xef, 0xa7,
	0x76, 0xed, 0x9b, 0xe3, 0x2f, 0x00, 0x00, 0xff, 0xff, 0xa0, 0xa6, 0x08, 0x40, 0xf3, 0x00, 0x00,
	0x00,
}
