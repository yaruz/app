package sn_account

import (
	"encoding"

	"github.com/golang/protobuf/proto"
)

var _ encoding.BinaryMarshaler = (*SNAccount)(nil)
var _ encoding.BinaryUnmarshaler = (*SNAccount)(nil)

func (e *SNAccount) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}

func (e *SNAccount) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}
