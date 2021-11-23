package date_value

import "time"

const (
	EntityName = "dateValue"
	TableName  = "date_value"
)

// DateValue ...
type DateValue struct {
	ID         uint      `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityID   uint      `gorm:"type:bigint not null REFERENCES \"entity\"(id);index:idx_date_value,priority:1" json:"entityID"`
	PropertyID uint      `gorm:"type:bigint not null;index:idx_date_value,priority:2" json:"propertyID"`
	Value      time.Time `gorm:"type:date not null;index" json:"value"`
}

func (e *DateValue) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *DateValue {
	return &DateValue{}
}
