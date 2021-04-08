package property_group

import "time"

const (
	EntityName = "propertyGroup"
	TableName  = "property_group"
)

// PropertyGroup ...
type PropertyGroup struct {
	ID          uint       `gorm:"type:bigint;primaryKey" json:"id"`
	Sysname     string     `gorm:"type:varchar(100) not null;unique;index" json:"sysname"`
	TSouirceID  uint       `sql:"type:bigint not null REFERENCES \"t_source\"(id)" gorm:"index" json:"tSourceId"`
	SortOrder   uint       `gorm:"type:int not null default 9999" json:"sortOrder"`
	Name        string     `gorm:"-" json:"name"`
	Description string     `gorm:"-" json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *PropertyGroup) TableName() string {
	return TableName
}

// New func is a constructor for the PropertyGroup
func New() *PropertyGroup {
	return &PropertyGroup{}
}
