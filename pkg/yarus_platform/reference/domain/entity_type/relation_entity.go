package entity_type

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
)

const (
	RelationEntityName = "relation"
	RelationTableName  = "property"
)

// Relation ...
type Relation struct {
	property.Property
	UndependedEntityType *EntityType `gorm:"-" json:"undependedEntityType"`
	DependedEntityType   *EntityType `gorm:"-" json:"dependedEntityType"`
}

func (e *Relation) TableName() string {
	return RelationTableName
}

// NewRelation func is a constructor for the Property
func NewRelation() *Relation {
	return &Relation{
		Property: property.Property{
			PropertyTypeID: property_type.IDRelation,
		},
	}
}

func (e Relation) Validate() error {
	err := validation.ValidateStruct(&e,
		validation.Field(&e.Property.PropertyTypeID, validation.Required, validation.In(uint(property_type.IDRelation))),
		validation.Field(&e.Property.IsRange, validation.In(false)),
	)
	if err != nil {
		return err
	}
	return e.Property.Validate()
}

func (e *Relation) SetUndependedEntityType(typeID *EntityType) {
	e.UndependedEntityType = typeID
}

func (e *Relation) SetDependedEntityType(typeID *EntityType) {
	e.DependedEntityType = typeID
}

func (e *Relation) AfterFind() error {
	return e.Property.AfterFind()
}

func (e *Relation) BeforeSave() error {
	return e.Property.BeforeSave()
}
