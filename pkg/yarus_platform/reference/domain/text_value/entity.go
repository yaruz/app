package text_value

import (
	"time"

	"gorm.io/gorm"
)

const (
	EntityName = "textValue"
	TableName  = "text_value"
)

// TextValue ...
type TextValue struct {
	ID           uint `gorm:"type:bigserial;primaryKey" json:"id"`
	TextSourceID uint `gorm:"type:bigint not null REFERENCES \"text_source\"(id);uniqueIndex:un_text_value;index" json:"textSourceID"`
	TextLangID   uint `gorm:"type:integer ;uniqueIndex:un_text_value" json:"langID"` // not null REFERENCES \"text_lang\"(id) [syntax error (SQLSTATE 42601)]

	Value     string         `gorm:"type:text not null" json:"value"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *TextValue) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *TextValue {
	return &TextValue{}
}
