package entity

import (
	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"
)

type RelationValue struct {
	entity_type.Relation
	Value []uint
}

func (v *RelationValue) SetValue(value interface{}) error {
	if v.Property.PropertyTypeID == 0 {
		return errors.Errorf("Can not set value to RelationValue: Property does not set.")
	}

	entitiesIDs, ok := value.([]uint)
	if !ok {
		return errors.Errorf("Can not cast value of relation into a []uint. Value = %v.", value)
	}
	v.Value = entitiesIDs

	return nil
}
