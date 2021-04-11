package property_view_type

import "time"

const (
	EntityName = "propertyViewType"
	TableName  = "property_view_type"
)

// PropertyViewType ...
type PropertyViewType struct {
	ID          uint       `gorm:"type:bigint;primaryKey" json:"id"`
	Sysname     string     `gorm:"type:varchar(100) not null;unique;index" json:"sysname"`
	TSouirceID  uint       `sql:"type:bigint not null REFERENCES \"t_source\"(id)" gorm:"index" json:"tSourceId"`
	Name        string     `gorm:"-" json:"name"`
	Description string     `gorm:"-" json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *PropertyViewType) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *PropertyViewType {
	return &PropertyViewType{}
}