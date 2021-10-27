package user

import (
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
)

const (
	EntityType = "user"
)

// User is the user entity
type User struct {
	entity.Entity
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// New func is a constructor for the User
func New(entity entity.Entity) *User {
	return &User{
		Entity: entity,
	}
}

func (e User) Validate() error {
	return nil
}
