package int_value

const (
	EntityName = "intValue"
	TableName  = "int_value"
)

// IntValue ...
type IntValue struct {
	ID         uint `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityID   uint `gorm:"type:bigint not null REFERENCES \"entity\"(id)" json:"entityID"`
	PropertyID uint `gorm:"type:bigint not null" json:"propertyID"`
	Value      int  `gorm:"type:bigint not null" json:"value"`
}

func (e *IntValue) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *IntValue {
	return &IntValue{}
}
