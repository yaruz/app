package user

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	EntityName = "user"
	TableName  = "user"
)

// User is the user entity
type User struct {
	ID        uint       `gorm:"primaryKey"`
	Name      string     `gorm:"type:varchar(100) not null;unique;index" json:"username"`
	Passhash  string     `gorm:"type:bytea not null" json:"-"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (e User) TableName() string {
	return TableName
}

// New func is a constructor for the User
func New() *User {
	return &User{}
}

func (e User) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Name, validation.Required, validation.Length(2, 100), is.Alpha),
	)
}
