package user

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
)

const (
	EntityType           = "User"
	PropertySysnameEmail = "User.Email"
	PropertySysnamePhone = "User.Phone"
)

// User is the user entity
type User struct {
	*entity.Entity
	ID    uint
	Email string
	Phone uint
}

var _ entity.Searchable = (*User)(nil)

func (e User) EntityType() string {
	return EntityType
}

func (e User) Validate() error {
	return nil
}

func (e *User) GetMapNameSysname() map[string]string {
	return map[string]string{
		"Email": PropertySysnameEmail,
		"Phone": PropertySysnamePhone,
	}
}

func (e *User) SetEmail(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnameEmail, 1)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 1); err != nil {
		return err
	}

	e.Email = value
	return nil
}

func (e *User) SetPhone(ctx context.Context, value uint) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnamePhone, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.Phone = value
	return nil
}
