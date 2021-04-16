package property_group

import "time"

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
