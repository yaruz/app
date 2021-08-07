package property_group

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/sysname"
	"gorm.io/gorm"
)

const (
	EntityName = "propertyGroup"
	TableName  = "property_group"
)

// PropertyGroup ...
type PropertyGroup struct {
	ID                  uint           `gorm:"type:bigserial;primaryKey" json:"id"`
	Sysname             string         `gorm:"type:varchar(100) not null;unique;index" json:"sysname"`
	NameSourceID        *uint          `gorm:"type:bigint null REFERENCES \"text_source\"(id);index" json:"-"`
	DescriptionSourceID *uint          `gorm:"type:bigint null REFERENCES \"text_source\"(id);index" json:"-"`
	SortOrder           uint           `gorm:"type:smallint;not null;default:9999" json:"sortOrder"`
	Name                *string        `gorm:"-" json:"name"`
	Description         *string        `gorm:"-" json:"description"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
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
		validation.Field(&e.Sysname, sysname.ValidationRules...),
	)
}
