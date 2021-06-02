package entity_type_relation

import (
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
)

const (
	EntityName = "EntityTypeRelation"
	TableName  = "property"
)

// EntityTypeRelation ...
type EntityTypeRelation struct {
	property.Property
	RelatedEntityType entity_type.EntityType `json:"relatedEntityType"`
	IsDependence      bool                   `json:"isDependence"`
}

func (e *EntityTypeRelation) TableName() string {
	return TableName
}

// New func is a constructor for the Property
func New() *EntityTypeRelation {
	return &EntityTypeRelation{}
}

func (e EntityTypeRelation) Validate() error {
	return e.Property.Validate()
}

func (e *EntityTypeRelation) AfterFind() error {
	return e.Property.AfterFind()
}

func (e *EntityTypeRelation) BeforeSave() error {
	return e.Property.BeforeSave()
}
