package float_value

const (
	EntityName = "floatValue"
	TableName  = "float_value"
)

// FloatValue ...
type FloatValue struct {
	ID         uint    `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityID   uint    `gorm:"type:bigint not null REFERENCES \"entity\"(id);uniqueIndex:un_float_value" json:"entityID"`
	PropertyID uint    `gorm:"type:bigint not null;uniqueIndex:un_float_value" json:"propertyID"`
	Value      float64 `gorm:"type:double precision not null" json:"value"`
}

func (e *FloatValue) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *FloatValue {
	return &FloatValue{}
}
