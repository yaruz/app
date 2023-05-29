package tg_session

import (
	"encoding"

	"github.com/golang/protobuf/proto"
)

var _ encoding.BinaryMarshaler = (*TgSession)(nil)
var _ encoding.BinaryUnmarshaler = (*TgSession)(nil)

func (e *TgSession) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}

func (e *TgSession) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}
