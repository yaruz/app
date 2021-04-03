package proto

import (
	"encoding"

	"github.com/golang/protobuf/proto"
)

var _ encoding.BinaryMarshaler = (*User)(nil)
var _ encoding.BinaryUnmarshaler = (*User)(nil)

func (e *User) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}

func (e *User) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}
