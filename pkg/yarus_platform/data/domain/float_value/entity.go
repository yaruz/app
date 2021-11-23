package float_value

const (
	EntityName = "floatValue"
	TableName  = "float_value"
)

// FloatValue ...
type FloatValue struct {
	ID         uint    `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityID   uint    `gorm:"type:bigint not null REFERENCES \"entity\"(id);index:idx_float_value,priority:1" json:"entityID"`
	PropertyID uint    `gorm:"type:bigint not null;index:idx_float_value,priority:2" json:"propertyID"`
	Value      float64 `gorm:"type:double precision not null;index" json:"value"`
}

func (e *FloatValue) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *FloatValue {
	return &FloatValue{}
}
