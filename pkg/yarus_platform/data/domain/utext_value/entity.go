package utext_value

const (
	EntityName = "utextValue"
	TableName  = "utext_value"
)

// UTextValue ...
type UTextValue struct {
	ID         uint   `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityID   uint   `gorm:"type:bigint not null REFERENCES \"entity\"(id);index:idx_utext_value,priority:1" json:"entityID"`
	PropertyID uint   `gorm:"type:bigint not null;index:idx_utext_value,priority:2" json:"propertyID"`
	Value      string `gorm:"type:text not null;index" json:"value"`
}

func (e *UTextValue) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *UTextValue {
	return &UTextValue{}
}
