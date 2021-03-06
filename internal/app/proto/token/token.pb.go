// Code generated by protoc-gen-go. DO NOT EDIT.
// source: token/token.proto

package token

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

type Token struct {
	AccessToken          string               `protobuf:"bytes,1,opt,name=AccessToken,proto3" json:"AccessToken,omitempty"`
	TokenType            string               `protobuf:"bytes,2,opt,name=TokenType,proto3" json:"TokenType,omitempty"`
	RefreshToken         string               `protobuf:"bytes,3,opt,name=RefreshToken,proto3" json:"RefreshToken,omitempty"`
	Expiry               *timestamp.Timestamp `protobuf:"bytes,4,opt,name=Expiry,proto3" json:"Expiry,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Token) Reset()         { *m = Token{} }
func (m *Token) String() string { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()    {}
func (*Token) Descriptor() ([]byte, []int) {
	return fileDescriptor_6e2ef433bb3fdc80, []int{0}
}

func (m *Token) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Token.Unmarshal(m, b)
}
func (m *Token) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Token.Marshal(b, m, deterministic)
}
func (m *Token) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Token.Merge(m, src)
}
func (m *Token) XXX_Size() int {
	return xxx_messageInfo_Token.Size(m)
}
func (m *Token) XXX_DiscardUnknown() {
	xxx_messageInfo_Token.DiscardUnknown(m)
}

var xxx_messageInfo_Token proto.InternalMessageInfo

func (m *Token) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

func (m *Token) GetTokenType() string {
	if m != nil {
		return m.TokenType
	}
	return ""
}

func (m *Token) GetRefreshToken() string {
	if m != nil {
		return m.RefreshToken
	}
	return ""
}

func (m *Token) GetExpiry() *timestamp.Timestamp {
	if m != nil {
		return m.Expiry
	}
	return nil
}

func init() {
	proto.RegisterType((*Token)(nil), "token.Token")
}

func init() {
	proto.RegisterFile("token/token.proto", fileDescriptor_6e2ef433bb3fdc80)
}

var fileDescriptor_6e2ef433bb3fdc80 = []byte{
	// 207 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0xc9, 0xcf, 0x4e,
	0xcd, 0xd3, 0x07, 0x93, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0xac, 0x60, 0x8e, 0x94, 0x7c,
	0x7a, 0x7e, 0x7e, 0x7a, 0x4e, 0xaa, 0x3e, 0x58, 0x30, 0xa9, 0x34, 0x4d, 0xbf, 0x24, 0x33, 0x37,
	0xb5, 0xb8, 0x24, 0x31, 0xb7, 0x00, 0xa2, 0x4e, 0x69, 0x3e, 0x23, 0x17, 0x6b, 0x08, 0x48, 0xa9,
	0x90, 0x02, 0x17, 0xb7, 0x63, 0x72, 0x72, 0x6a, 0x71, 0x31, 0x98, 0x2b, 0xc1, 0xa8, 0xc0, 0xa8,
	0xc1, 0x19, 0x84, 0x2c, 0x24, 0x24, 0xc3, 0xc5, 0x09, 0x66, 0x84, 0x54, 0x16, 0xa4, 0x4a, 0x30,
	0x81, 0xe5, 0x11, 0x02, 0x42, 0x4a, 0x5c, 0x3c, 0x41, 0xa9, 0x69, 0x45, 0xa9, 0xc5, 0x19, 0x10,
	0x03, 0x98, 0xc1, 0x0a, 0x50, 0xc4, 0x84, 0x8c, 0xb8, 0xd8, 0x5c, 0x2b, 0x0a, 0x32, 0x8b, 0x2a,
	0x25, 0x58, 0x14, 0x18, 0x35, 0xb8, 0x8d, 0xa4, 0xf4, 0x20, 0xee, 0xd3, 0x83, 0xb9, 0x4f, 0x2f,
	0x04, 0xe6, 0xbe, 0x20, 0xa8, 0x4a, 0x27, 0xfd, 0x28, 0xdd, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24,
	0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0xca, 0xc4, 0xa2, 0xd2, 0x2a, 0xfd, 0xc4, 0x82, 0x02, 0xfd, 0xcc,
	0xbc, 0x92, 0xd4, 0xa2, 0xbc, 0xc4, 0x1c, 0x30, 0x07, 0xac, 0x1f, 0x12, 0x00, 0x49, 0x6c, 0x60,
	0x8e, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x9e, 0xb3, 0xab, 0x8e, 0x16, 0x01, 0x00, 0x00,
}
