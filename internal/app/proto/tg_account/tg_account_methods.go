package tg_account

import (
	"encoding"

	"github.com/golang/protobuf/proto"
)

var _ encoding.BinaryMarshaler = (*TgAccount)(nil)
var _ encoding.BinaryUnmarshaler = (*TgAccount)(nil)
var _ encoding.BinaryMarshaler = (*AuthSession)(nil)
var _ encoding.BinaryUnmarshaler = (*AuthSession)(nil)

func (e *TgAccount) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}

func (e *TgAccount) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}

func (e *AuthSession) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}

func (e *AuthSession) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}
