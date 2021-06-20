package entity

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

const (
	EntityName = "entity"
	TableName  = "entity"
)

// Entity ...
type Entity struct {
	ID           uint           `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityTypeID uint           `gorm:"type:bigint not null;index" json:"entityTypeID"`
	PropertiesB  datatypes.JSON `json:"-"`
	//PropertiesValuesMap map[uint]interface{}   `gorm:"-" json:"propertiesValuesMap"`
	PropertiesValues map[string]PropertyValues  `gorm:"-" json:"propertiesValues"`
	RelationsValues  map[string]RelatedEntities `gorm:"-" json:"relationsValues"`
	CreatedAt        time.Time                  `json:"createdAt"`
	UpdatedAt        time.Time                  `json:"updatedAt"`
	DeletedAt        *time.Time                 `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *Entity) TableName() string {
	return TableName
}

// New func is a constructor for the Entity
func New() *Entity {
	return &Entity{}
}

func (e *Entity) AfterFind() error {
	return e.propertiesB2PropertiesValues()
}

func (e *Entity) BeforeSave() error {
	return e.propertiesValues2PropertiesB()
}

func (e *Entity) propertiesB2PropertiesValues() error {
	jsonb, err := e.PropertiesB.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonb, &e.PropertiesValues)
}

func (e *Entity) propertiesValues2PropertiesB() error {
	jsonb, err := json.Marshal(&e.PropertiesValues)
	if err != nil {
		return err
	}
	return e.PropertiesB.UnmarshalJSON(jsonb)
}
