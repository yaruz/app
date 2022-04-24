package proto

import (
	"encoding"

	"github.com/golang/protobuf/proto"
)

var _ encoding.BinaryMarshaler = (*SnAccount)(nil)
var _ encoding.BinaryUnmarshaler = (*SnAccount)(nil)

func (e *SnAccount) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}

func (e *SnAccount) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}
