package relation

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
)

const (
	EntityName = "relation"
	TableName  = "property"
)

// Relation ...
type Relation struct {
	property.Property
	UndependedEntityType *entity_type.EntityType `gorm:"many2many:entity_type;" json:"entityType"`
	DependedEntityType   *entity_type.EntityType `gorm:"many2many:entity_type;" json:"dependedEntityType"`
}

func (e *Relation) TableName() string {
	return TableName
}

// New func is a constructor for the Property
func New() *Relation {
	return &Relation{}
}

func (e Relation) Validate() error {
	err := validation.ValidateStruct(&e,
		validation.Field(&e.PropertyTypeID, validation.Required, validation.In(property_type.IDRelation)),
		validation.Field(&e.IsRange, validation.In(false)),
	)
	if err != nil {
		return err
	}
	return e.Property.Validate()
}

func (e *Relation) AfterFind() error {
	return e.Property.AfterFind()
}

func (e *Relation) BeforeSave() error {
	return e.Property.BeforeSave()
}
