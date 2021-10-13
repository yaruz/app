package entity

import (
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/bool_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/date_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/float_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/int_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/time_value"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"

	"gorm.io/datatypes"
)

const (
	EntityName = "entity"
	TableName  = "entity"
)

// Entity ...
// Значения свойст:
// 	пишем в PropertiesValuesMap
// 	храним в PropertiesB (автоматом конвертим туда-сюда)
// 	PropertiesValues - красивый и удобный список со свойствами, только для чтения, инициализируется в сервисе entity.Service.EntityInit()
// 	RelationsValues - красивый и удобный список со связями - // - // -
//
type Entity struct {
	ID                  uint                     `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityTypeID        uint                     `gorm:"type:bigint not null;index" json:"entityTypeID"`
	PropertiesB         datatypes.JSON           `gorm:"index:idx_properties_b,type:gin" json:"-"`
	PropertiesValuesMap map[uint]interface{}     `gorm:"-" json:"-"`
	PropertiesValues    map[uint]PropertyValue   `gorm:"-" json:"propertiesValues"`
	RelationsValues     map[uint]RelationValue   `gorm:"-" json:"relationsValues"`
	TextValues          []text_value.TextValue   `json:"-"`
	BoolValues          []bool_value.BoolValue   `json:"-"`
	IntValues           []int_value.IntValue     `json:"-"`
	FloatValues         []float_value.FloatValue `json:"-"`
	DateValues          []date_value.DateValue   `json:"-"`
	TimeValues          []time_value.TimeValue   `json:"-"`
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

func (e *Entity) SetPropertyValue(propertyValue *PropertyValue) {
	propertyID := propertyValue.Property.ID

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

func (e *Entity) SetRelationValue(relationValue *RelationValue) {
	propertyID := relationValue.Property.ID

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

// Удаляет как значения свойств, так и значения связей
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

func (e *Entity) deletePropertyBoolValues(propertyID uint) {
	for i := range e.BoolValues {
		if e.BoolValues[i].PropertyID == propertyID {
			e.BoolValues = append(e.BoolValues[:i], e.BoolValues[i+1:]...)
		}
	}
}

func (e *Entity) deletePropertyIntValues(propertyID uint) {
	for i := range e.IntValues {
		if e.IntValues[i].PropertyID == propertyID {
			e.IntValues = append(e.IntValues[:i], e.IntValues[i+1:]...)
		}
	}
}

func (e *Entity) deletePropertyFloatValues(propertyID uint) {
	for i := range e.FloatValues {
		if e.FloatValues[i].PropertyID == propertyID {
			e.FloatValues = append(e.FloatValues[:i], e.FloatValues[i+1:]...)
		}
	}
}

func (e *Entity) deletePropertyDateValues(propertyID uint) {
	for i := range e.DateValues {
		if e.DateValues[i].PropertyID == propertyID {
			e.DateValues = append(e.DateValues[:i], e.DateValues[i+1:]...)
		}
	}
}

func (e *Entity) deletePropertyTimeValues(propertyID uint) {
	for i := range e.TimeValues {
		if e.TimeValues[i].PropertyID == propertyID {
			e.TimeValues = append(e.TimeValues[:i], e.TimeValues[i+1:]...)
		}
	}
}

func (e *Entity) deletePropertyTextValues(propertyID uint) {
	for i := range e.TextValues {
		if e.TextValues[i].PropertyID == propertyID {
			e.TextValues = append(e.TextValues[:i], e.TextValues[i+1:]...)
		}
	}
}
