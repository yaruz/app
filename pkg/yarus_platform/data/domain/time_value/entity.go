package time_value

import "time"

const (
	EntityName = "timeValue"
	TableName  = "time_value"
)

// TimeValue ...
type TimeValue struct {
	ID         uint      `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityID   uint      `gorm:"type:bigint not null REFERENCES \"entity\"(id);uniqueIndex:un_time_value" json:"entityID"`
	PropertyID uint      `gorm:"type:bigint not null;uniqueIndex:un_time_value" json:"propertyID"`
	Value      time.Time `gorm:"type:timestamptz not null" json:"value"`
}

func (e *TimeValue) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *TimeValue {
	return &TimeValue{}
}
