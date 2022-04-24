package proto

import (
	"encoding"

	"github.com/golang/protobuf/proto"
)

var _ encoding.BinaryMarshaler = (*JwtClaims)(nil)
var _ encoding.BinaryUnmarshaler = (*JwtClaims)(nil)

var _ encoding.BinaryMarshaler = (*Account)(nil)
var _ encoding.BinaryUnmarshaler = (*Account)(nil)

var _ encoding.BinaryMarshaler = (*AccountSettings)(nil)
var _ encoding.BinaryUnmarshaler = (*AccountSettings)(nil)

func (e *AccountSettings) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}

func (e *AccountSettings) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}

func (e *JwtClaims) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}

func (e *JwtClaims) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}

func (e *Account) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}

func (e *Account) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}
