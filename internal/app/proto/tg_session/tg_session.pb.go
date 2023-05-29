// Code generated by protoc-gen-go. DO NOT EDIT.
// source: tg_session/tg_session.proto

package tg_session

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type TgSession struct {
	IsAuthorized         bool     `protobuf:"varint,1,opt,name=IsAuthorized,proto3" json:"IsAuthorized,omitempty"`
	Session              []byte   `protobuf:"bytes,2,opt,name=Session,proto3" json:"Session,omitempty"`
	ID                   string   `protobuf:"bytes,3,opt,name=ID,proto3" json:"ID,omitempty"`
	Phone                string   `protobuf:"bytes,4,opt,name=Phone,proto3" json:"Phone,omitempty"`
	PhoneCodeHash        string   `protobuf:"bytes,5,opt,name=PhoneCodeHash,proto3" json:"PhoneCodeHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TgSession) Reset()         { *m = TgSession{} }
func (m *TgSession) String() string { return proto.CompactTextString(m) }
func (*TgSession) ProtoMessage()    {}
func (*TgSession) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b16f07807b864a3, []int{0}
}

func (m *TgSession) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TgSession.Unmarshal(m, b)
}
func (m *TgSession) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TgSession.Marshal(b, m, deterministic)
}
func (m *TgSession) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TgSession.Merge(m, src)
}
func (m *TgSession) XXX_Size() int {
	return xxx_messageInfo_TgSession.Size(m)
}
func (m *TgSession) XXX_DiscardUnknown() {
	xxx_messageInfo_TgSession.DiscardUnknown(m)
}

var xxx_messageInfo_TgSession proto.InternalMessageInfo

func (m *TgSession) GetIsAuthorized() bool {
	if m != nil {
		return m.IsAuthorized
	}
	return false
}

func (m *TgSession) GetSession() []byte {
	if m != nil {
		return m.Session
	}
	return nil
}

func (m *TgSession) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *TgSession) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *TgSession) GetPhoneCodeHash() string {
	if m != nil {
		return m.PhoneCodeHash
	}
	return ""
}

func init() {
	proto.RegisterType((*TgSession)(nil), "tg_session.TgSession")
}

func init() {
	proto.RegisterFile("tg_session/tg_session.proto", fileDescriptor_7b16f07807b864a3)
}

var fileDescriptor_7b16f07807b864a3 = []byte{
	// 193 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2e, 0x49, 0x8f, 0x2f,
	0x4e, 0x2d, 0x2e, 0xce, 0xcc, 0xcf, 0xd3, 0x47, 0x30, 0xf5, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85,
	0xb8, 0x10, 0x22, 0x4a, 0x53, 0x19, 0xb9, 0x38, 0x43, 0xd2, 0x83, 0x21, 0x3c, 0x21, 0x25, 0x2e,
	0x1e, 0xcf, 0x62, 0xc7, 0xd2, 0x92, 0x8c, 0xfc, 0xa2, 0xcc, 0xaa, 0xd4, 0x14, 0x09, 0x46, 0x05,
	0x46, 0x0d, 0x8e, 0x20, 0x14, 0x31, 0x21, 0x09, 0x2e, 0x76, 0xa8, 0x72, 0x09, 0x26, 0x05, 0x46,
	0x0d, 0x9e, 0x20, 0x18, 0x57, 0x88, 0x8f, 0x8b, 0xc9, 0xd3, 0x45, 0x82, 0x59, 0x81, 0x51, 0x83,
	0x33, 0x88, 0xc9, 0xd3, 0x45, 0x48, 0x84, 0x8b, 0x35, 0x20, 0x23, 0x3f, 0x2f, 0x55, 0x82, 0x05,
	0x2c, 0x04, 0xe1, 0x08, 0xa9, 0x70, 0xf1, 0x82, 0x19, 0xce, 0xf9, 0x29, 0xa9, 0x1e, 0x89, 0xc5,
	0x19, 0x12, 0xac, 0x60, 0x59, 0x54, 0x41, 0x27, 0x93, 0x28, 0xa3, 0xf4, 0xcc, 0x92, 0x8c, 0xd2,
	0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0xca, 0xc4, 0xa2, 0xd2, 0x2a, 0xfd, 0xc4, 0x82, 0x02, 0xfd,
	0xcc, 0xbc, 0x92, 0xd4, 0xa2, 0xbc, 0xc4, 0x1c, 0x30, 0x07, 0xec, 0x1d, 0x24, 0xff, 0x25, 0xb1,
	0x81, 0x45, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x36, 0x73, 0xf0, 0xef, 0xff, 0x00, 0x00,
	0x00,
}
