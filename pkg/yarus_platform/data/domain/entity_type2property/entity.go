package entity_type2property

import "time"

const (
	EntityName = "entityType2Property"
	TableName  = "entity_type2property"
)

// EntityType2Property ...
type EntityType2Property struct {
	EntityTypeID uint       `sql:"type:bigint not null REFERENCES \"entity_type\"(id);primary_key" json:"entityTypeId"`
	PropertyID   uint       `sql:"type:bigint not null REFERENCES \"property\"(id);primary_key" json:"propertyId"`
	IsComposite  bool       `gorm:"type:boolean not null;" json:"isComposite"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *EntityType2Property) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *EntityType2Property {
	return &EntityType2Property{}
}
