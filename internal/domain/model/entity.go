package model

const (
	EntityName = "model"
	TableName  = "car_model"
)

// Post is the user entity
type Model struct {
	ID      uint   `gorm:"column:id_car_model" json:"id"`
	MarkID  uint   `gorm:"column:id_car_mark" json:"markId"`
	Name    string `gorm:"type:varchar(255)" json:"name"`
	NameRus string `gorm:"type:varchar(255)" json:"nameRus"`
	TypeID  uint   `gorm:"column:id_car_type" json:"typeId"`
}

func (e Model) TableName() string {
	return TableName
}

// New func is a constructor for the Post
func New() *Model {
	return &Model{}
}
