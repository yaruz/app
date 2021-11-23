package text_value

const (
	EntityName = "textValue"
	TableName  = "text_value"
)

// TextValue ...
type TextValue struct {
	ID            uint   `gorm:"type:bigserial;primaryKey" json:"id"`
	EntityID      uint   `gorm:"type:bigint not null REFERENCES \"entity\"(id);index:idx_text_value,priority:1;index:idx_text_value_inst,priority:1" json:"entityID"`
	LangID        uint   `gorm:"type:smallint not null;index:idx_text_value,priority:3;index:idx_text_value_inst,priority:2" json:"langID"`
	PropertyID    uint   `gorm:"type:bigint not null;index:idx_text_value,priority:2" json:"propertyID"`
	Value         string `gorm:"type:text not null" json:"value"`
	ValueTsvector string `gorm:"type:tsvector not null;index:textsearch_idx,type:GIN" json:"-"`
}

func (e *TextValue) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *TextValue {
	return &TextValue{}
}
