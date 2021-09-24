package entity

import (
	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
)

type RelationValue struct {
	*entity_type.Relation
	Value []uint
}

func newRelationValue(relation *entity_type.Relation, value interface{}) (*RelationValue, error) {
	relationValue := &RelationValue{
		Relation: relation,
	}

	if err := relationValue.SetValue(value); err != nil {
		return nil, err
	}
	return relationValue, nil
}

func (v *RelationValue) SetValue(value interface{}) (err error) {
	if v.Property.PropertyTypeID == 0 {
		return errors.Errorf("Can not set value to RelationValue: Property does not set.")
	}

	if v.Property.PropertyTypeID != property_type.IDRelation {
		return errors.Errorf("Can not set value to RelationValue: PropertyTypeID does not equal IDRelation.")
	}

	v.Value, err = property.GetRelationValue(value)

	return err
}
