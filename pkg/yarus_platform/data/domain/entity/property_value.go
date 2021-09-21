package entity

import (
	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
)

type PropertyValue struct {
	property.Property
	Value interface{} //	<simple type> || []<simple type>
}

func (v *PropertyValue) SetValue(value interface{}) error {
	if v.Property.PropertyTypeID == 0 {
		return errors.Errorf("Can not set value to PropertyValue: Property does not set.")
	}

	switch v.Property.PropertyTypeID {
	case property_type.IDBoolean:
		val, ok := value.(bool)
		if !ok {
			return errors.Errorf("Can not cast value of PropertyValue to bool. Value = %v.", value)
		}
		v.Value = val
	case property_type.IDInt:
		val, ok := value.(int)
		if !ok {
			return errors.Errorf("Can not cast value of PropertyValue to int. Value = %v.", value)
		}
		v.Value = val
	case property_type.IDFloat:
		val, ok := value.(float64)
		if !ok {
			return errors.Errorf("Can not cast value of PropertyValue to float64. Value = %v.", value)
		}
		v.Value = val
	case property_type.IDDate:
	case property_type.IDTimestamp:
	case property_type.IDText:
		val, ok := value.(string)
		if !ok {
			return errors.Errorf("Can not cast value of PropertyValue to string. Value = %v.", value)
		}
		v.Value = val
	default:
		return errors.Errorf("Can not set value to PropertyValue: unknown PropertyTypeID = %v.", v.Property.PropertyTypeID)
	}

	return nil
}
