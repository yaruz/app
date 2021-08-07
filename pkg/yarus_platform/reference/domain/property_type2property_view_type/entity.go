package property_type2property_view_type

import (
	"time"

	"gorm.io/gorm"
)

const (
	EntityName = "propertyType2PropertyViewType"
	TableName  = "property_type2property_view_type"
)

// PropertyType2PropertyViewType ...
type PropertyType2PropertyViewType struct {
	PropertyTypeID     uint           `gorm:"type:bigint not null REFERENCES \"property_type\"(id);primaryKey" json:"propertyTypeID"`
	PropertyViewTypeID uint           `gorm:"type:bigint not null REFERENCES \"property_view_type\"(id);primaryKey" json:"propertyViewTypeID"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *PropertyType2PropertyViewType) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *PropertyType2PropertyViewType {
	return &PropertyType2PropertyViewType{}
}
