package text_lang

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gorm.io/gorm"
)

const (
	EntityName = "textLang"
	TableName  = "text_lang"
)

// TextLang ...
type TextLang struct {
	ID        uint           `gorm:"type:smallserial;primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(100) not null;unique;index" json:"code"`
	Name      string         `gorm:"type:varchar(100) not null" json:"name"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *TextLang) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *TextLang {
	return &TextLang{}
}

func (e TextLang) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Code, validation.Length(0, 3), is.Alpha, is.LowerCase),
	)
}
