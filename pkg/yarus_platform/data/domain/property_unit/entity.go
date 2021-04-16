package property_unit

import "time"

const (
	EntityName = "propertyUnit"
	TableName  = "property_unit"
)

// PropertyUnit ...
type PropertyUnit struct {
	ID                  uint       `gorm:"type:bigint;primaryKey" json:"id"`
	Sysname             string     `gorm:"type:varchar(100) not null;unique;index" json:"sysname"`
	NameSourceID        uint       `sql:"type:bigint not null REFERENCES \"text_source\"(id)" gorm:"index" json:"-"`
	DescriptionSourceID uint       `sql:"type:bigint not null REFERENCES \"text_source\"(id)" gorm:"index" json:"-"`
	Name                string     `gorm:"-" json:"name"`
	Description         string     `gorm:"-" json:"description"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
	DeletedAt           *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *PropertyUnit) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *PropertyUnit {
	return &PropertyUnit{}
}
