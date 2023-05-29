package jwt

import (
	"encoding"

	"github.com/golang/protobuf/proto"
)

var _ encoding.BinaryMarshaler = (*Claims)(nil)
var _ encoding.BinaryUnmarshaler = (*Claims)(nil)

var _ encoding.BinaryMarshaler = (*RegisteredClaims)(nil)
var _ encoding.BinaryUnmarshaler = (*RegisteredClaims)(nil)

var _ encoding.BinaryMarshaler = (*TokenData)(nil)
var _ encoding.BinaryUnmarshaler = (*TokenData)(nil)

func (e *Claims) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}
func (e *Claims) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}

func (e *RegisteredClaims) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}
func (e *RegisteredClaims) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}

func (e *TokenData) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}
func (e *TokenData) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}
