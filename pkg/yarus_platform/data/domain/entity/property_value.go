package entity

import (
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
)

type PropertyValue struct {
	property.Property
	Value interface{} //	<simple type> || []<simple type>
}
