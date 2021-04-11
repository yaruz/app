package t_text

import (
	"time"
)

const (
	EntityName = "tText"
	TableName  = "t_text"
)

// TText ...
type TText struct {
	ID         uint       `gorm:"type:bigint;primaryKey" json:"id"`
	TSourceID  uint       `sql:"type:bigint not null REFERENCES \"t_source\"(id)" gorm:"index:un_t_text,unique" json:"tSourceId"`
	LangID     uint       `gorm:"type:smallint not null;index:un_t_text,unique" json:"langID"`
	PropertyID uint       `sql:"type:bigint REFERENCES \"property\"(id)" gorm:"index:un_t_text,unique" json:"propertyID"`
	Value      string     `gorm:"type:text not null" json:"value"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *TText) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *TText {
	return &TText{}
}
