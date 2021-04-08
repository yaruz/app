package t_string

import (
	"time"
)

const (
	EntityName = "tString"
	TableName  = "t_string"
)

// TString ...
type TString struct {
	ID         uint       `gorm:"type:bigint;primaryKey" json:"id"`
	TSourceID  uint       `sql:"type:bigint not null REFERENCES \"t_source\"(id)" gorm:"index:un_t_string,unique" json:"tSourceId"`
	LangID     uint       `gorm:"type:smallint not null;index:un_t_string,unique" json:"langID"`
	PropertyID uint       `sql:"type:bigint REFERENCES \"property\"(id)" gorm:"index:un_t_string,unique" json:"propertyID"`
	Value      string     `gorm:"type:varchar(255) not null" json:"value"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *TString) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *TString {
	return &TString{}
}
