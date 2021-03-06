package property_unit

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/sysname"
	"gorm.io/gorm"
)

const (
	EntityName  = "propertyUnit"
	TableName   = "property_unit"
	SysnameItem = "item"
)

// PropertyUnit ...
type PropertyUnit struct {
	ID                  uint           `gorm:"type:bigserial;primaryKey" json:"id"`
	Sysname             string         `gorm:"type:varchar(100) not null;unique;index" json:"sysname"`
	NameSourceID        *uint          `gorm:"type:bigint null REFERENCES \"text_source\"(id);index" json:"-"`
	DescriptionSourceID *uint          `gorm:"type:bigint null REFERENCES \"text_source\"(id);index" json:"-"`
	Name                *string        `gorm:"-" json:"name"`
	Description         *string        `gorm:"-" json:"description"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *PropertyUnit) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *PropertyUnit {
	return &PropertyUnit{}
}

func (e PropertyUnit) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Sysname, sysname.ValidationRules...),
	)
}
