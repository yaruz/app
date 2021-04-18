package text_value

import (
	"time"
)

const (
	EntityName = "textValue"
	TableName  = "text_value"
)

// TextValue ...
type TextValue struct {
	ID           uint       `gorm:"type:bigint;primaryKey" json:"id"`
	TextSourceID uint       `sql:"type:bigint not null REFERENCES \"text_source\"(id);unique_index:un_text_value" gorm:"index" json:"textSourceID"`
	LangID       uint       `sql:"type:smallint not null;unique_index:un_text_value" json:"langID"`
	PropertyID   uint       `sql:"type:bigint REFERENCES \"property\"(id)" gorm:"index" json:"propertyID"`
	Value        string     `gorm:"type:text not null" json:"value"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *TextValue) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *TextValue {
	return &TextValue{}
}
