// Code generated by protoc-gen-go. DO NOT EDIT.
// source: session/session.proto

package session

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	account "github.com/yaruz/app/internal/app/proto/account"
	jwt "github.com/yaruz/app/internal/app/proto/jwt"
	user "github.com/yaruz/app/internal/app/proto/user"
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

type Session struct {
	ID                   string                   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	AccountSettings      *account.AccountSettings `protobuf:"bytes,2,opt,name=AccountSettings,proto3" json:"AccountSettings,omitempty"`
	JwtClaims            *jwt.Claims              `protobuf:"bytes,3,opt,name=JwtClaims,proto3" json:"JwtClaims,omitempty"`
	User                 *user.User               `protobuf:"bytes,4,opt,name=User,proto3" json:"User,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *Session) Reset()         { *m = Session{} }
func (m *Session) String() string { return proto.CompactTextString(m) }
func (*Session) ProtoMessage()    {}
func (*Session) Descriptor() ([]byte, []int) {
	return fileDescriptor_e18c811121f5946b, []int{0}
}

func (m *Session) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Session.Unmarshal(m, b)
}
func (m *Session) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Session.Marshal(b, m, deterministic)
}
func (m *Session) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Session.Merge(m, src)
}
func (m *Session) XXX_Size() int {
	return xxx_messageInfo_Session.Size(m)
}
func (m *Session) XXX_DiscardUnknown() {
	xxx_messageInfo_Session.DiscardUnknown(m)
}

var xxx_messageInfo_Session proto.InternalMessageInfo

func (m *Session) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Session) GetAccountSettings() *account.AccountSettings {
	if m != nil {
		return m.AccountSettings
	}
	return nil
}

func (m *Session) GetJwtClaims() *jwt.Claims {
	if m != nil {
		return m.JwtClaims
	}
	return nil
}

func (m *Session) GetUser() *user.User {
	if m != nil {
		return m.User
	}
	return nil
}

func init() {
	proto.RegisterType((*Session)(nil), "session.Session")
}

func init() {
	proto.RegisterFile("session/session.proto", fileDescriptor_e18c811121f5946b)
}

var fileDescriptor_e18c811121f5946b = []byte{
	// 225 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2d, 0x4e, 0x2d, 0x2e,
	0xce, 0xcc, 0xcf, 0xd3, 0x87, 0xd2, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0xec, 0x50, 0xae,
	0x14, 0x7f, 0x69, 0x71, 0x6a, 0x91, 0x3e, 0x88, 0x80, 0xc8, 0x48, 0x89, 0x26, 0x26, 0x27, 0xe7,
	0x97, 0xe6, 0x95, 0xe8, 0x43, 0x69, 0xa8, 0x30, 0x6f, 0x56, 0x79, 0x89, 0x7e, 0x56, 0x39, 0x94,
	0xab, 0xb4, 0x82, 0x91, 0x8b, 0x3d, 0x18, 0x62, 0x84, 0x10, 0x1f, 0x17, 0x93, 0xa7, 0x8b, 0x04,
	0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x93, 0xa7, 0x8b, 0x90, 0x13, 0x17, 0xbf, 0x23, 0x44, 0x6f,
	0x70, 0x6a, 0x49, 0x49, 0x66, 0x5e, 0x7a, 0xb1, 0x04, 0x93, 0x02, 0xa3, 0x06, 0xb7, 0x91, 0x84,
	0x1e, 0xcc, 0x4c, 0x34, 0xf9, 0x20, 0x74, 0x0d, 0x42, 0x9a, 0x5c, 0x9c, 0x5e, 0xe5, 0x25, 0xce,
	0x39, 0x89, 0x99, 0xb9, 0xc5, 0x12, 0xcc, 0x60, 0xdd, 0xdc, 0x7a, 0x20, 0xeb, 0x21, 0x42, 0x41,
	0x08, 0x59, 0x21, 0x39, 0x2e, 0x96, 0xd0, 0xe2, 0xd4, 0x22, 0x09, 0x16, 0xb0, 0x2a, 0x2e, 0x3d,
	0xb0, 0x5f, 0x40, 0x22, 0x41, 0x60, 0x71, 0x27, 0xc3, 0x28, 0xfd, 0xf4, 0xcc, 0x92, 0x8c, 0xd2,
	0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0xca, 0xc4, 0xa2, 0xd2, 0x2a, 0xfd, 0xc4, 0x82, 0x02, 0xfd,
	0xcc, 0xbc, 0x92, 0xd4, 0xa2, 0xbc, 0xc4, 0x1c, 0x30, 0x07, 0xec, 0x2b, 0x58, 0x18, 0x25, 0xb1,
	0x81, 0xb9, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x05, 0xa5, 0x6f, 0x14, 0x3d, 0x01, 0x00,
	0x00,
}
