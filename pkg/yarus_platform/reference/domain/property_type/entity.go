package property_type

import (
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_view_type"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	EntityName       = "propertyType"
	TableName        = "property_type"
	SysnameRelation  = "relation"
	SysnameBoolean   = "boolean"
	SysnameInt       = "int"
	SysnameBigint    = "bigint"
	SysnameFloat     = "float"
	SysnameDate      = "date"
	SysnameTimestamp = "timestamp"
	SysnameText      = "text"
	IDRelation       = 1
	IDBoolean        = 2
	IDInt            = 3
	IDBigint         = 4
	IDFloat          = 5
	IDDate           = 6
	IDTimestamp      = 7
	IDText           = 8
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
		validation.Field(&e.Sysname, domain.SysnameValidationRules...),
	)
}

func (e PropertyType) GetPropertyViewTypes() ([]property_view_type.PropertyViewType, error) {
	return e.PropertyViewTypes, nil
}

// PropertyType2PropertyViewType ...
//type PropertyType2PropertyViewType struct {
//	PropertyTypeID     uint `gorm:"type:bigint not null REFERENCES \"property_type\"(id);primaryKey" json:"propertyTypeID"`
//	PropertyViewTypeID uint `gorm:"type:bigint not null REFERENCES \"property_view_type\"(id);primaryKey" json:"propertyViewTypeID"`
//}
