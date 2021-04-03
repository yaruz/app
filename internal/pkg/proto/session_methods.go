package proto

import (
	"encoding"

	"github.com/golang/protobuf/proto"
)

var _ encoding.BinaryMarshaler = (*Session)(nil)
var _ encoding.BinaryUnmarshaler = (*Session)(nil)

func (e *Session) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}

func (e *Session) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}
