package text_value

const (
	EntityName = "textValue"
	TableName  = "text_value"
)

// TextValue ...
type TextValue struct {
	ID         uint   `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityID   uint   `gorm:"type:bigint not null REFERENCES \"entity\"(id);uniqueIndex:un_text_value" json:"entityID"`
	LangID     uint   `gorm:"type:smallint not null;uniqueIndex:un_text_value" json:"langID"`
	PropertyID uint   `gorm:"type:bigint not null;uniqueIndex:un_text_value" json:"propertyID"`
	Value      string `gorm:"type:text not null" json:"value"`
}

func (e *TextValue) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *TextValue {
	return &TextValue{}
}
