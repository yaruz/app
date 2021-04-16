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
	ID           uint                 `gorm:"type:bigint;primaryKey" json:"id"`
	EntityTypeID uint                 `sql:"type:bigint not null REFERENCES \"entity_type\"(id)" gorm:"index" json:"entityTypeID"`
	PropertiesB  datatypes.JSON       `json:"-"`
	Properties   map[uint]interface{} `gorm:"-" json:"properties"`
	CreatedAt    time.Time            `json:"createdAt"`
	UpdatedAt    time.Time            `json:"updatedAt"`
	DeletedAt    *time.Time           `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *Entity) TableName() string {
	return TableName
}

// New func is a constructor for the Entity
func New() *Entity {
	return &Entity{}
}

func (e *Entity) Init() error {
	return e.propertiesB2Properties()
}

func (e *Entity) Prepare4Save() error {
	return e.properties2PropertiesB()
}

func (e *Entity) propertiesB2Properties() error {
	jsonb, err := e.PropertiesB.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonb, &e.Properties)
}

func (e *Entity) properties2PropertiesB() error {
	jsonb, err := json.Marshal(&e.Properties)
	if err != nil {
		return err
	}
	return e.PropertiesB.UnmarshalJSON(jsonb)
}
