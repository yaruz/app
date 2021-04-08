package t_source

import (
	"time"
)

const (
	EntityName = "tSource"
	TableName  = "t_source"
)

// TSource ...
type TSource struct {
	ID          uint       `gorm:"type:bigint;primaryKey" json:"id"`
	SouirceID   uint       `sql:"type:bigint not null" gorm:"index:un_t_source__source_id_table_field,unique" json:"souirceID"`
	Table       string     `gorm:"type:varchar(50) not null;index:un_t_source__source_id_table_field,unique" json:"table"`
	Field       string     `gorm:"type:varchar(50) not null;index:un_t_source__source_id_table_field,unique" json:"field"`
	Name        string     `gorm:"-" json:"name"`
	Description string     `gorm:"-" json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *TSource) TableName() string {
	return TableName
}

// New func is a constructor for the EntityType
func New() *TSource {
	return &TSource{}
}
