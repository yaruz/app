package entity

import (
	"sort"

	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
)

type RelationValue struct {
	*entity_type.Relation
	Value []uint
}

func newRelationValue(relation *entity_type.Relation, value []uint) (*RelationValue, error) {
	relationValue := &RelationValue{
		Relation: relation,
	}

	if err := relationValue.SetValue(value); err != nil {
		return nil, err
	}
	return relationValue, nil
}

func (v *RelationValue) propertyTypeIDCheck() error {
	if v.Property.PropertyTypeID == 0 {
		return errors.Errorf("Can not set value to RelationValue: Property does not set.")
	}

	if v.Property.PropertyTypeID != property_type.IDRelation {
		return errors.Errorf("Can not set value to RelationValue: PropertyTypeID does not equal IDRelation.")
	}
	return nil
}

func (v *RelationValue) SetValueByInterface(value interface{}) (err error) {
	if err = v.propertyTypeIDCheck(); err != nil {
		return err
	}
	v.Value, err = property.GetRelationValue(value)
	return err
}

func (v *RelationValue) SetValue(value []uint) error {
	if err := v.propertyTypeIDCheck(); err != nil {
		return err
	}
	v.Value = value
	return nil
}

func (v *RelationValue) BindRelatedEntity(entityID uint) {
	if v.Value == nil {
		v.Value = make([]uint, 0, 1)
	}
	v.Value = append(v.Value, entityID)
}

func (v *RelationValue) BindRelatedEntities(entityIDs []uint) {
	if v.Value == nil {
		v.Value = make([]uint, 0, len(entityIDs))
	}
	v.Value = append(v.Value, entityIDs...)
}

func (v *RelationValue) UnbindRelatedEntity(entityID uint) error {
	if v.Value == nil {
		return yaruserror.ErrNotFound
	}
	if _, ok := v.Value[entityID]; !ok {
		return yaruserror.ErrNotFound
	}

}

func (v *RelationValue) UnbindRelatedEntities(entityIDs []uint) {
}

func (v *RelationValue) sortValue() {
	sort.Slice(v.Value, func(i, j int) bool { return v.Value[i] < v.Value[j] })
}

func (v *RelationValue) searchValue(value uint) (int, bool) {
	intSlice := make([]int, len(v.Value))

	for i, val := range v.Value {
		intSlice[i] = int(val)
	}

	indx := sort.SearchInts(intSlice, int(value))
	return indx, len(v.Value) > indx && v.Value[indx] == value
}
