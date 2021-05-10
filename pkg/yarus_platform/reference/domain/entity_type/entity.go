package entity_type

import (
	"regexp"
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	EntityName = "entityType"
	TableName  = "entity_type"
)

// EntityType ...
type EntityType struct {
	ID                  uint                `gorm:"type:bigserial;primaryKey" json:"id"`
	Sysname             string              `gorm:"type:varchar(100) not null;unique;index" json:"sysname"`
	NameSourceID        *uint               `gorm:"type:bigint null REFERENCES \"text_source\"(id)" gorm:"index" json:"-"`
	DescriptionSourceID *uint               `gorm:"type:bigint null REFERENCES \"text_source\"(id);index" json:"-"`
	Name                *string             `gorm:"-" json:"name"`
	Description         *string             `gorm:"-" json:"description"`
	Properties          []property.Property `gorm:"many2many:entity_type2property;"`
	CreatedAt           time.Time           `json:"createdAt"`
	UpdatedAt           time.Time           `json:"updatedAt"`
	DeletedAt           *time.Time          `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *EntityType) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *EntityType {
	return &EntityType{}
}

func (e EntityType) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Sysname, validation.Required, validation.Length(2, 100), validation.Match(regexp.MustCompile("^[a-z0-9_]+$"))),
	)
}
