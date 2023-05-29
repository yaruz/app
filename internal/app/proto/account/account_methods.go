package account

import (
	"encoding"

	"github.com/golang/protobuf/proto"
)

var _ encoding.BinaryMarshaler = (*AccountSettings)(nil)
var _ encoding.BinaryUnmarshaler = (*AccountSettings)(nil)

func (e *AccountSettings) MarshalBinary() (data []byte, err error) {
	return proto.Marshal(e)
}

func (e *AccountSettings) UnmarshalBinary(data []byte) (err error) {
	return proto.Unmarshal(data, e)
}
