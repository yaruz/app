package date_value

import "time"

const (
	EntityName = "dateValue"
	TableName  = "date_value"
)

// DateValue ...
type DateValue struct {
	ID         uint      `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityID   uint      `gorm:"type:bigint not null REFERENCES \"entity\"(id)" json:"entityID"`
	PropertyID uint      `gorm:"type:bigint not null" json:"propertyID"`
	Value      time.Time `gorm:"type:date not null" json:"value"`
}

func (e *DateValue) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *DateValue {
	return &DateValue{}
}
