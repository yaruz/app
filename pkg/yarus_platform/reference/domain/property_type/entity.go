package property_type

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_view_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/sysname"
	"gorm.io/gorm"
)

const (
	EntityName      = "propertyType"
	TableName       = "property_type"
	SysnameRelation = "relation"
	SysnameBoolean  = "boolean"
	SysnameInt      = "int"
	SysnameFloat    = "float"
	SysnameDate     = "date"
	SysnameTime     = "time"
	SysnameText     = "text"
	IDRelation      = 1
	IDBoolean       = 2
	IDInt           = 3
	IDFloat         = 4
	IDDate          = 5
	IDTime          = 6
	IDText          = 7
	Count           = 6 // для табличной интерпретации
)

var ValidTypesIDs = []interface{}{
	uint(IDRelation),
	uint(IDBoolean),
	uint(IDInt),
	uint(IDFloat),
	uint(IDDate),
	uint(IDTime),
	uint(IDText),
}

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
	DeletedAt           gorm.DeletedAt                        `gorm:"index" json:"deletedAt,omitempty"`
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
		validation.Field(&e.Sysname, sysname.ValidationRules...),
	)
}

func (e PropertyType) GetPropertyViewTypes() ([]property_view_type.PropertyViewType, error) {
	return e.PropertyViewTypes, nil
}
