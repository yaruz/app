package tg_account

import (
	"encoding"

	"github.com/golang/protobuf/proto"
)

var _ encoding.BinaryMarshaler = (*TgAccount)(nil)
var _ encoding.BinaryUnmarshaler = (*TgAccount)(nil)

func (e *TgAccount) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}

func (e *TgAccount) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}
