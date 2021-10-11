package entity

import (
	"encoding/json"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/bool_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/date_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/float_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/int_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/time_value"

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

func (e *Entity) AfterFind() error {
	return e.propertiesB2PropertiesValuesMap()
}

func (e *Entity) BeforeSave() error {
	return e.propertiesValuesMap2PropertiesB()
}

func (e *Entity) propertiesB2PropertiesValuesMap() error {
	jsonb, err := e.PropertiesB.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonb, &e.PropertiesValuesMap)
}

func (e *Entity) propertiesValuesMap2PropertiesB() error {
	jsonb, err := json.Marshal(&e.PropertiesValuesMap)
	if err != nil {
		return err
	}
	return e.PropertiesB.UnmarshalJSON(jsonb)
}

func (e *Entity) SetPropertyValue(propertyValue *PropertyValue) {
	propertyID := propertyValue.Property.ID

	if e.PropertiesValuesMap == nil {
		e.PropertiesValuesMap = make(map[uint]interface{}, 1)
	}
	if e.PropertiesValues == nil {
		e.PropertiesValues = make(map[uint]PropertyValue, 1)
	}

	e.PropertiesValuesMap[propertyID] = propertyValue.Value
	e.PropertiesValues[propertyID] = *propertyValue
}

func (e *Entity) SetRelationValue(relationValue *RelationValue) {
	propertyID := relationValue.Property.ID

	if e.PropertiesValuesMap == nil {
		e.PropertiesValuesMap = make(map[uint]interface{}, 1)
	}
	if e.RelationsValues == nil {
		e.RelationsValues = make(map[uint]RelationValue, 1)
	}

	e.PropertiesValuesMap[propertyID] = relationValue.Value
	e.RelationsValues[propertyID] = *relationValue
}

// Удаляет как значения свойств, так и значения связей
func (e *Entity) DeletePropertyValue(propertyID uint) {

	if _, ok := e.PropertiesValuesMap[propertyID]; ok {
		delete(e.PropertiesValuesMap, propertyID)
	}

	if _, ok := e.PropertiesValues[propertyID]; ok {
		delete(e.PropertiesValues, propertyID)
	}

	if _, ok := e.RelationsValues[propertyID]; ok {
		delete(e.RelationsValues, propertyID)
	}
}
