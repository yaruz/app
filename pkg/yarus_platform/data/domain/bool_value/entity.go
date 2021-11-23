package bool_value

const (
	EntityName = "boolValue"
	TableName  = "bool_value"
)

// BoolValue ...
type BoolValue struct {
	ID         uint `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityID   uint `gorm:"type:bigint not null REFERENCES \"entity\"(id);index:idx_bool_value,priority:1" json:"entityID"`
	PropertyID uint `gorm:"type:bigint not null;index:idx_bool_value,priority:2" json:"propertyID"`
	Value      bool `gorm:"type:boolean not null" json:"value"`
}

func (e *BoolValue) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *BoolValue {
	return &BoolValue{}
}
