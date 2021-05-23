package property_type

import (
	"regexp"
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_view_type"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	EntityName       = "propertyType"
	TableName        = "property_type"
	SysnameBoolean   = "boolean"
	SysnameInt       = "int"
	SysnameBigint    = "bigint"
	SysnameFloat     = "float"
	SysnameDate      = "date"
	SysnameTimestamp = "timestamp"
	SysnameText      = "text"
)

var fl float64

// PropertyType ...
type PropertyType struct {
	ID                  uint                                  `gorm:"type:bigserial;primaryKey" json:"id"`
	Sysname             string                                `gorm:"type:varchar(100) not null;unique;index" json:"sysname"`
	NameSourceID        *uint                                 `gorm:"type:bigint null REFERENCES \"text_source\"(id);index" json:"-"`
	DescriptionSourceID *uint                                 `gorm:"type:bigint null REFERENCES \"text_source\"(id);index" json:"-"`
	Name                *string                               `gorm:"-" json:"name"`
	Description         *string                               `gorm:"-" json:"description"`
	PropertyViewTypes   []property_view_type.PropertyViewType `gorm:"many2many:property_type2property_view_type;"`
	CreatedAt           time.Time                             `json:"createdAt"`
	UpdatedAt           time.Time                             `json:"updatedAt"`
	DeletedAt           *time.Time                            `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *PropertyType) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *PropertyType {
	return &PropertyType{}
}

func (e PropertyType) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Sysname, validation.Required, validation.Length(2, 100), validation.Match(regexp.MustCompile("^[a-z0-9_]+$"))),
	)
}

func (e PropertyType) GetPropertyViewTypes() ([]property_view_type.PropertyViewType, error) {
	return e.PropertyViewTypes, nil
}

// PropertyType2PropertyViewType ...
type PropertyType2PropertyViewType struct {
	PropertyTypeID     uint `gorm:"type:bigint not null REFERENCES \"property_type\"(id);primaryKey" json:"propertyTypeID"`
	PropertyViewTypeID uint `gorm:"type:bigint not null REFERENCES \"property_view_type\"(id);primaryKey" json:"propertyViewTypeID"`
}
