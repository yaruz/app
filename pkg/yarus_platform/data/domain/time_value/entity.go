package time_value

import "time"

const (
	EntityName = "timeValue"
	TableName  = "time_value"
)

// TimeValue ...
type TimeValue struct {
	ID         uint      `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityID   uint      `gorm:"type:bigint not null REFERENCES \"entity\"(id);index:idx_time_value,priority:1" json:"entityID"`
	PropertyID uint      `gorm:"type:bigint not null;index:idx_time_value,priority:2" json:"propertyID"`
	Value      time.Time `gorm:"type:timestamptz not null;index" json:"value"`
}

func (e *TimeValue) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *TimeValue {
	return &TimeValue{}
}
