package entity

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	EntityName = "entity"
	TableName  = "entity"
)

// Entity ...
type Entity struct {
	ID                  uint                   `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityTypeID        uint                   `gorm:"type:bigint not null;index" json:"entityTypeID"`
	PropertiesB         datatypes.JSON         `json:"-"`
	PropertiesValuesMap map[uint]interface{}   `gorm:"-" json:"-"`
	PropertiesValues    map[uint]PropertyValue `gorm:"-" json:"propertiesValues"`
	RelationsValues     map[uint]RelationValue `gorm:"-" json:"relationsValues"`
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
