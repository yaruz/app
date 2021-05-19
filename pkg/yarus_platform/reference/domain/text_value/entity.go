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
	TextSourceID uint       `gorm:"type:bigint not null REFERENCES \"text_source\"(id);uniqueIndex:un_text_value;index" json:"textSourceID"`
	LangID       uint       `gorm:"type:smallint not null;uniqueIndex:un_text_value" json:"langID"`
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