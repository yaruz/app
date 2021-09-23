package entity

import (
	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
)

type PropertyValue struct {
	*property.Property
	Value interface{} //	<simple type> || []<simple type>
}

func (v *PropertyValue) SetValue(value interface{}, langID uint) (err error) {
	if v.Property.PropertyTypeID == 0 {
		return errors.Errorf("Can not set value to PropertyValue: Property does not set.")
	}

	switch v.Property.PropertyTypeID {
	case property_type.IDBoolean:
		v.Value, err = property.GetValueBool(value)
	case property_type.IDInt:
		v.Value, err = property.GetValueInt(value)
	case property_type.IDFloat:
		v.Value, err = property.GetValueFloat(value)
	case property_type.IDDate:
		v.Value, err = property.GetValueDate(value)
	case property_type.IDTime:
		v.Value, err = property.GetValueTime(value)
	case property_type.IDText:
		v.Value, err = property.GetValueText(value)
	default:
		err = errors.Errorf("Can not set value to PropertyValue: unknown PropertyTypeID = %v.", v.Property.PropertyTypeID)
	}

	return err
}
