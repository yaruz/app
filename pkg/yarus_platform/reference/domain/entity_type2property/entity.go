package entity_type2property

import (
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"
)

const (
	EntityName = "entityType2Property"
	TableName  = "entity_type2property"
)

// EntityType2Property ...
type EntityType2Property struct {
	EntityTypeID uint                    `gorm:"type:bigint not null REFERENCES \"entity_type\"(id);primaryKey" json:"entityTypeId"`
	PropertyID   uint                    `gorm:"type:bigint not null REFERENCES \"property\"(id);primaryKey" json:"propertyId"`
	IsDependent  bool                    `gorm:"type:boolean not null default false;primaryKey;" json:"isDependent"`
	EntityType   *entity_type.EntityType `json:"entityType"`
	CreatedAt    time.Time               `json:"createdAt"`
}

func (e *EntityType2Property) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *EntityType2Property {
	return &EntityType2Property{}
}
