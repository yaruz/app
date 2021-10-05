package entity

import (
	"encoding/json"
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

	if err := relationValue.propertyTypeIDCheck(); err != nil {
		return nil, err
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

	if v.Relation.Property.PropertyTypeID != property_type.IDRelation {
		return errors.Errorf("Can not set value to RelationValue: PropertyTypeID does not equal IDRelation.")
	}
	return nil
}

func (v *RelationValue) SetValueByInterface(value interface{}) (err error) {
	valueUint, err := property.GetRelationValue(value)
	if err != nil {
		return err
	}
	return v.SetValue(valueUint)
}

func (v *RelationValue) SetValue(value []uint) error {
	v.Value = value
	v.sortValue()
	return nil
}

func (v *RelationValue) AddValue(value uint) error {
	if v.Value == nil {
		v.Value = make([]uint, 0, 1)
	}

	i, ok := v.SearchValue(value)
	if ok {
		return yaruserror.ErrAlreadyExists
	}

	if i == 0 {
		v.Value = append([]uint{value}, v.Value...)
	} else if i < v.Len() {
		v.Value = append(v.Value[:i+1], v.Value[i:]...)
		v.Value[i] = value
	} else {
		v.Value = append(v.Value, value)
	}
	return nil
}

func (v *RelationValue) AddValues(values []uint, isStopIfErrAlreadyExists bool) error {
	if values == nil || len(values) == 0 {
		return yaruserror.ErrEmptyParams
	}

	if v.Value == nil {
		v.Value = make([]uint, 0, len(values))
	}
	alreadyExists := make(map[int]uint)

	for i, value := range values {
		if _, ok := v.SearchValue(value); ok {
			alreadyExists[i] = value
		}
	}

	if len(alreadyExists) == 0 || !isStopIfErrAlreadyExists {
		for i, id := range values {
			if _, ok := alreadyExists[i]; ok {
				continue
			}
			if err := v.AddValue(id); err != nil {
				if !errors.Is(err, yaruserror.ErrAlreadyExists) {
					return err
				}
				alreadyExists[i] = id
			}
		}
	}

	if len(alreadyExists) > 0 {
		alreadyExistsB, err := json.Marshal(alreadyExists)
		if err != nil {
			return err
		}
		return errors.Wrap(yaruserror.ErrAlreadyExists, string(alreadyExistsB))
	}
	return nil
}

func (v *RelationValue) RemoveValue(Value uint) error {
	i, ok := v.SearchValue(Value)
	if !ok {
		return yaruserror.ErrNotFound
	}

	v.Value = append(v.Value[:i], v.Value[i+1:]...)
	return nil
}

func (v *RelationValue) RemoveValues(entityIDs []uint, isStopIfErrNotFound bool) error {
	if entityIDs == nil || len(entityIDs) == 0 {
		return yaruserror.ErrEmptyParams
	}

	notFound := make(map[int]uint)
	for i, id := range entityIDs {
		if _, ok := v.SearchValue(id); !ok {
			notFound[i] = id
		}
	}

	if len(notFound) == 0 || !isStopIfErrNotFound {
		for i, id := range entityIDs {
			if _, ok := notFound[i]; ok {
				continue
			}
			if err := v.RemoveValue(id); err != nil {
				return err
			}
		}
	}

	if len(notFound) > 0 {
		notFoundB, err := json.Marshal(notFound)
		if err != nil {
			return err
		}
		return errors.Wrap(yaruserror.ErrNotFound, string(notFoundB))
	}
	return nil
}

func (v *RelationValue) sortValue() {
	v.CheckValueNotNil()
	sort.Slice(v.Value, func(i, j int) bool { return v.Value[i] < v.Value[j] })
}

func (v *RelationValue) SearchValue(value uint) (int, bool) {
	v.CheckValueNotNil()
	intSlice := make([]int, len(v.Value))

	for i, val := range v.Value {
		intSlice[i] = int(val)
	}

	indx := sort.SearchInts(intSlice, int(value))
	return indx, len(v.Value) > indx && v.Value[indx] == value
}

func (v *RelationValue) CheckValueNotNil() {
	if v.Value == nil {
		v.Value = make([]uint, 0)
	}
}

func (v *RelationValue) Len() int {
	v.CheckValueNotNil()
	return len(v.Value)
}
