package entity

import (
	"encoding/json"
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"

	"gorm.io/datatypes"
	"gorm.io/gorm"
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
	ID                  uint                   `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityTypeID        uint                   `gorm:"type:bigint not null;index" json:"entityTypeID"`
	PropertiesB         datatypes.JSON         `json:"-"`
	PropertiesValuesMap map[uint]interface{}   `gorm:"-" json:"-"`
	PropertiesValues    map[uint]PropertyValue `gorm:"-" json:"propertiesValues"`
	RelationsValues     map[uint]RelationValue `gorm:"-" json:"relationsValues"`
	TextValues          []text_value.TextValue `json:"-"`
	CreatedAt           time.Time              `json:"createdAt"`
	UpdatedAt           time.Time              `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt         `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *Entity) TableName() string {
	return TableName
}

// New func is a constructor for the Entity
func New() *Entity {
	return &Entity{}
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

func (e *Entity) setPropertyValue(propertyValue *PropertyValue) {
	propertyID := propertyValue.Property.ID

	e.PropertiesValuesMap[propertyID] = propertyValue.Value
	e.PropertiesValues[propertyID] = *propertyValue
}

func (e *Entity) setRelationValue(relationValue *RelationValue) {
	propertyID := relationValue.Property.ID

	e.PropertiesValuesMap[propertyID] = relationValue.Value
	e.RelationsValues[propertyID] = *relationValue
}

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
