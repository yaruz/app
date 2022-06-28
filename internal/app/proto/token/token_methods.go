package token

import (
	"encoding"

	"github.com/golang/protobuf/proto"
)

var _ encoding.BinaryMarshaler = (*Token)(nil)
var _ encoding.BinaryUnmarshaler = (*Token)(nil)

func (e *Token) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}

func (e *Token) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}
