package text_source

import (
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_value"
	"gorm.io/gorm"
)

const (
	EntityName = "textSource"
	TableName  = "text_source"
)

// TextSource ...
type TextSource struct {
	ID         uint                   `gorm:"type:bigserial;primaryKey" json:"id"`
	TextValue  *text_value.TextValue  `json:"textValue"`
	TextValues []text_value.TextValue `json:"textValues"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *TextSource) TableName() string {
	return TableName
}

// New func is a constructor for the TextSource
func New() *TextSource {
	return &TextSource{}
}
