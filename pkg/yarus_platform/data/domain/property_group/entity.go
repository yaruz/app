package property_group

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	EntityName = "propertyGroup"
	TableName  = "property_group"
)

// PropertyGroup ...
type PropertyGroup struct {
	ID                  uint       `gorm:"type:bigint;primaryKey" json:"id"`
	Sysname             string     `gorm:"type:varchar(100) not null;unique;index" json:"sysname"`
	NameSourceID        uint       `sql:"type:bigint not null REFERENCES \"text_source\"(id)" gorm:"index" json:"-"`
	DescriptionSourceID uint       `sql:"type:bigint not null REFERENCES \"text_source\"(id)" gorm:"index" json:"-"`
	SortOrder           uint       `gorm:"type:smallint not null default 9999" json:"sortOrder"`
	Name                string     `gorm:"-" json:"name"`
	Description         string     `gorm:"-" json:"description"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
	DeletedAt           *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *PropertyGroup) TableName() string {
	return TableName
}

// New func is a constructor for the PropertyGroup
func New() *PropertyGroup {
	return &PropertyGroup{}
}

func (e PropertyGroup) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Sysname, validation.Required, validation.Length(2, 100), validation.Match(regexp.MustCompile("^[a-z0-9_]+$"))),
	)
}
