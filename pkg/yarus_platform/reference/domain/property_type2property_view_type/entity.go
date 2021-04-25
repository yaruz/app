package property_type2property_view_type

import "time"

const (
	EntityName = "propertyType2PropertyViewType"
	TableName  = "property_type2property_view_type"
)

// PropertyType2PropertyViewType ...
type PropertyType2PropertyViewType struct {
	PropertyTypeID     uint       `sql:"type:bigint not null REFERENCES \"property_type\"(id);primary_key" json:"propertyTypeID"`
	PropertyViewTypeID uint       `sql:"type:bigint not null REFERENCES \"property_view_type\"(id);primary_key" json:"propertyViewTypeID"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	DeletedAt          *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *PropertyType2PropertyViewType) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *PropertyType2PropertyViewType {
	return &PropertyType2PropertyViewType{}
}
