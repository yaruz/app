package entity

import (
	"github.com/pkg/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/bool_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/date_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/float_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/int_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/time_value"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
)

const (
	EntityName = "entity"
	TableName  = "entity"
)

// Entity ...
// Значения свойст:
// 	пишем в соотв. слайс из ...Values
// 	храним в соотв. таблице ..._value
// 	PropertiesValues - красивый и удобный список со свойствами, только для чтения, инициализируется в сервисе entity.Service.EntityInit()
// 	RelationsValues - красивый и удобный список со связями - // - // -
//
type Entity struct {
	ID               uint                     `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityTypeID     uint                     `gorm:"type:bigint not null;index" json:"entityTypeID"`
	PropertiesValues map[uint]PropertyValue   `gorm:"-" json:"propertiesValues"`
	RelationsValues  map[uint]RelationValue   `gorm:"-" json:"relationsValues"`
	TextValues       []text_value.TextValue   `json:"-"`
	BoolValues       []bool_value.BoolValue   `json:"-"`
	IntValues        []int_value.IntValue     `json:"-"`
	FloatValues      []float_value.FloatValue `json:"-"`
	DateValues       []date_value.DateValue   `json:"-"`
	TimeValues       []time_value.TimeValue   `json:"-"`
}

func (e *Entity) TableName() string {
	return TableName
}

// New func is a constructor for the Entity
func New() *Entity {
	return &Entity{}
}

func (e Entity) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.EntityTypeID, validation.Required),
	)
}

// Присваивает PropertyValue
func (e *Entity) SetPropertyValue(propertyValue *PropertyValue) {
	propertyID := propertyValue.Property.ID
	e.DeletePropertyValues(propertyID)

	if e.PropertiesValues == nil {
		e.PropertiesValues = make(map[uint]PropertyValue, 1)
	}

	e.PropertiesValues[propertyID] = *propertyValue

	switch propertyValue.Property.PropertyTypeID {
	case property_type.IDBoolean:
		value, _ := property.GetValueBool(propertyValue.Value)
		e.BoolValues = append(e.BoolValues, bool_value.BoolValue{
			EntityID:   e.ID,
			PropertyID: propertyID,
			Value:      value,
		})
	case property_type.IDInt:
		value, _ := property.GetValueInt(propertyValue.Value)
		e.IntValues = append(e.IntValues, int_value.IntValue{
			EntityID:   e.ID,
			PropertyID: propertyID,
			Value:      value,
		})
	case property_type.IDFloat:
		value, _ := property.GetValueFloat(propertyValue.Value)
		e.FloatValues = append(e.FloatValues, float_value.FloatValue{
			EntityID:   e.ID,
			PropertyID: propertyID,
			Value:      value,
		})
	case property_type.IDDate:
		value, _ := property.GetValueDate(propertyValue.Value)
		e.DateValues = append(e.DateValues, date_value.DateValue{
			EntityID:   e.ID,
			PropertyID: propertyID,
			Value:      value,
		})
	case property_type.IDTime:
		value, _ := property.GetValueTime(propertyValue.Value)
		e.TimeValues = append(e.TimeValues, time_value.TimeValue{
			EntityID:   e.ID,
			PropertyID: propertyID,
			Value:      value,
		})
	case property_type.IDText:
		value, _ := property.GetValueText(propertyValue.Value)
		e.TextValues = append(e.TextValues, text_value.TextValue{
			EntityID:   e.ID,
			PropertyID: propertyID,
			Value:      value,
		})
	}
}

// Присваивает RelationValue
func (e *Entity) SetRelationValue(relationValue *RelationValue) {
	propertyID := relationValue.Property.ID
	e.DeletePropertyValues(propertyID)

	if e.RelationsValues == nil {
		e.RelationsValues = make(map[uint]RelationValue, 1)
	}

	e.RelationsValues[propertyID] = *relationValue

	for _, val := range relationValue.Value {
		value, _ := property.GetValueInt(val)
		e.IntValues = append(e.IntValues, int_value.IntValue{
			EntityID:   e.ID,
			PropertyID: propertyID,
			Value:      value,
		})
	}
}

// Удаляет все как значения свойств, так и значения связей для заданного propertyID
func (e *Entity) DeletePropertyValues(propertyID uint) {

	propertyValue, propOk := e.PropertiesValues[propertyID]
	_, relOk := e.RelationsValues[propertyID]

	if propOk {
		delete(e.PropertiesValues, propertyID)

		switch propertyValue.Property.PropertyTypeID {
		case property_type.IDBoolean:
			e.deletePropertyBoolValues(propertyID)
		case property_type.IDInt:
			e.deletePropertyIntValues(propertyID)
		case property_type.IDFloat:
			e.deletePropertyFloatValues(propertyID)
		case property_type.IDDate:
			e.deletePropertyDateValues(propertyID)
		case property_type.IDTime:
			e.deletePropertyTimeValues(propertyID)
		case property_type.IDText:
			e.deletePropertyTextValues(propertyID)
		}
	} else if relOk {
		delete(e.RelationsValues, propertyID)
		e.deletePropertyIntValues(propertyID)
	}
}

// Удаляет все значения BoolValues для заданного propertyID
func (e *Entity) deletePropertyBoolValues(propertyID uint) {
	for i := range e.BoolValues {
		if e.BoolValues[i].PropertyID == propertyID {
			e.BoolValues = append(e.BoolValues[:i], e.BoolValues[i+1:]...)
		}
	}
}

// Удаляет все значения IntValues для заданного propertyID
func (e *Entity) deletePropertyIntValues(propertyID uint) {
	for i := range e.IntValues {
		if e.IntValues[i].PropertyID == propertyID {
			e.IntValues = append(e.IntValues[:i], e.IntValues[i+1:]...)
		}
	}
}

// Удаляет все значения FloatValues для заданного propertyID
func (e *Entity) deletePropertyFloatValues(propertyID uint) {
	for i := range e.FloatValues {
		if e.FloatValues[i].PropertyID == propertyID {
			e.FloatValues = append(e.FloatValues[:i], e.FloatValues[i+1:]...)
		}
	}
}

// Удаляет все значения DateValues для заданного propertyID
func (e *Entity) deletePropertyDateValues(propertyID uint) {
	for i := range e.DateValues {
		if e.DateValues[i].PropertyID == propertyID {
			e.DateValues = append(e.DateValues[:i], e.DateValues[i+1:]...)
		}
	}
}

// Удаляет все значения TimeValues для заданного propertyID
func (e *Entity) deletePropertyTimeValues(propertyID uint) {
	for i := range e.TimeValues {
		if e.TimeValues[i].PropertyID == propertyID {
			e.TimeValues = append(e.TimeValues[:i], e.TimeValues[i+1:]...)
		}
	}
}

// Удаляет все значения TextValues для заданного propertyID
func (e *Entity) deletePropertyTextValues(propertyID uint) {
	for i := range e.TextValues {
		if e.TextValues[i].PropertyID == propertyID {
			e.TextValues = append(e.TextValues[:i], e.TextValues[i+1:]...)
		}
	}
}

// Возвращает количество всех значений всех свойств
func (e *Entity) PropertiesValuesCount() int {
	return len(e.BoolValues) + len(e.IntValues) + len(e.FloatValues) + len(e.DateValues) + len(e.TimeValues) + len(e.TextValues)
}

// Возвращает слайс ID всех свойств, значения которых заданны для entity
func (e *Entity) GetPropertiesIDs() []interface{} {
	propertiesIdsMap := make(map[uint]struct{}, e.PropertiesValuesCount())

	for _, val := range e.BoolValues {
		propertiesIdsMap[val.PropertyID] = struct{}{}
	}

	propertiesIds := make([]interface{}, 0, len(propertiesIdsMap))
	for id := range propertiesIdsMap {
		propertiesIds = append(propertiesIds, id)
	}

	return propertiesIds
}

// Возвращает слайс ID значений связей для заданного propertyID
func (e *Entity) GetRelationValues(propertyID uint) []uint {
	values := make([]uint, 0, 1)

	for _, intValue := range e.IntValues {
		if intValue.PropertyID == propertyID {
			values = append(values, uint(intValue.Value))
		}
	}
	return values
}

// Инициализация PropertiesValues и RelationsValues из всех слайсов значений свойств ...Values
// Запускаем после запуска Service.propertiesInit()
func (e *Entity) propertiesValuesInit(langID uint) error {

	for _, propertyValue := range e.BoolValues {
		if err := e.setValueToExistingPropertyValue(propertyValue.PropertyID, propertyValue.Value, langID); err != nil {
			return err
		}
	}

	for _, propertyValue := range e.IntValues {
		if err := e.setValueToExistingPropertyValue(propertyValue.PropertyID, propertyValue.Value, langID); err != nil {
			return err
		}
	}

	for _, propertyValue := range e.FloatValues {
		if err := e.setValueToExistingPropertyValue(propertyValue.PropertyID, propertyValue.Value, langID); err != nil {
			return err
		}
	}

	for _, propertyValue := range e.DateValues {
		if err := e.setValueToExistingPropertyValue(propertyValue.PropertyID, propertyValue.Value, langID); err != nil {
			return err
		}
	}

	for _, propertyValue := range e.TimeValues {
		if err := e.setValueToExistingPropertyValue(propertyValue.PropertyID, propertyValue.Value, langID); err != nil {
			return err
		}
	}

	for _, propertyValue := range e.TextValues {
		if err := e.setValueToExistingPropertyValue(propertyValue.PropertyID, propertyValue.Value, langID); err != nil {
			return err
		}
	}

	return nil
}

// Только для метода propertiesValuesInit()
// Уже к существующим после иннициализации экземплярам RelationValue и PropertyValue присваивает значение value
func (e *Entity) setValueToExistingPropertyValue(propertyID uint, value interface{}, langID uint) error {
	propertyValue, propOk := e.PropertiesValues[propertyID]
	relationValue, relOk := e.RelationsValues[propertyID]

	switch {
	case relOk:
		if err := relationValue.SetValueByInterface(value); err != nil {
			return errors.Wrapf(err, "Can not set value to PropertyValue. Property ID = %v; Value = %v.", propertyID, value)
		}
		e.RelationsValues[propertyID] = relationValue
	case propOk:
		if err := propertyValue.SetValue(value, langID); err != nil {
			return errors.Wrapf(err, "Can not set value to PropertyValue. Property ID = %v; Value = %v.", propertyID, value)
		}
		e.PropertiesValues[propertyID] = propertyValue
	default:
		return errors.Errorf("Property ID = %v not found.", propertyID)
	}
	return nil
}

// Создаёт и присваивает новые значения PropertyValue
// value - значение, не ссылка
func (e *Entity) SetValueForProperty(property *property.Property, value interface{}, langID uint) error {
	propertyValue, err := newPropertyValue(property, value, langID)
	if err != nil {
		return err
	}
	e.SetPropertyValue(propertyValue)
	return nil
}

// Создаёт и присваивает новые значения RelationValue
// value - значение, не ссылка, []uint - IDs связанных сущностей
func (e *Entity) SetValueForRelation(relation *entity_type.Relation, value []uint) error {
	relationValue, err := newRelationValue(relation, value)
	if err != nil {
		return err
	}
	e.SetRelationValue(relationValue)
	return nil
}

// По заданному значению relation привязываем relatedEntityID
func (e *Entity) BindRelatedEntityID(relation *entity_type.Relation, relatedEntityID uint) error {
	propertyID := relation.Property.ID

	if len(e.GetRelationValues(propertyID)) == 0 {
		return e.SetValueForRelation(relation, []uint{relatedEntityID})
	}

	relationsValues, ok := e.RelationsValues[propertyID]
	if !ok {
		return errors.Wrapf(yaruserror.ErrNotSet, "RelationsValues fo relation ID = %v not init", propertyID)
	}

	if err := relationsValues.AddValue(relatedEntityID); err != nil {
		return err
	}

	e.SetRelationValue(&relationsValues)
	return nil
}

// По заданному значению relation привязываем relatedEntityIDs
func (e *Entity) BindRelatedEntityIDs(relation *entity_type.Relation, relatedEntityIDs []uint) error {
	propertyID := relation.Property.ID

	if len(e.GetRelationValues(propertyID)) == 0 {
		return e.SetValueForRelation(relation, relatedEntityIDs)
	}

	relationsValues, ok := e.RelationsValues[propertyID]
	if !ok {
		return errors.Wrapf(yaruserror.ErrNotSet, "RelationsValues fo relation ID = %v not init", propertyID)
	}

	if err := relationsValues.AddValues(relatedEntityIDs, false); err != nil {
		return err
	}

	e.SetRelationValue(&relationsValues)
	return nil
}

// По заданному значению relation отвязываем relatedEntityID
func (e *Entity) UnbindRelatedEntityID(relation *entity_type.Relation, relatedEntityID uint) error {
	propertyID := relation.Property.ID

	if len(e.GetRelationValues(propertyID)) == 0 {
		return errors.Wrapf(yaruserror.ErrNotFound, "related entity with ID = %v not found", relatedEntityID)
	}

	relationsValues, ok := e.RelationsValues[propertyID]
	if !ok {
		return errors.Wrapf(yaruserror.ErrNotSet, "RelationsValues fo relation ID = %v not init", propertyID)
	}

	if err := relationsValues.RemoveValue(relatedEntityID); err != nil {
		return err
	}

	if len(relationsValues.Value) == 0 {
		e.DeletePropertyValues(propertyID)
	} else {
		e.SetRelationValue(&relationsValues)
	}

	return nil
}

// По заданному значению relation отвязываем relatedEntityIDs
func (e *Entity) UnbindRelatedEntityIDs(relation *entity_type.Relation, relatedEntityIDs []uint) error {
	propertyID := relation.Property.ID

	if len(e.GetRelationValues(propertyID)) == 0 {
		return errors.Wrapf(yaruserror.ErrNotFound, "related entity with ID = %v not found", relatedEntityIDs)
	}

	relationsValues, ok := e.RelationsValues[propertyID]
	if !ok {
		return errors.Wrapf(yaruserror.ErrNotSet, "RelationsValues fo relation ID = %v not init", propertyID)
	}

	return relationsValues.RemoveValues(relatedEntityIDs, false)
}
