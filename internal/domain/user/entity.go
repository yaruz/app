package user

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
)

const (
	EntityType               = "User"
	PropertySysnameEmail     = "User.Email"
	PropertySysnameAccountID = "User.AccountID"
)

// User is the user entity
type User struct {
	*entity.Entity
	ID        uint
	AccountID string
	Email     string
}

type AccountSettings struct {
	LangID uint
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
		"Email":     PropertySysnameEmail,
		"AccountID": PropertySysnameAccountID,
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

func (e *User) SetAccountID(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnameAccountID, 1)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 1); err != nil {
		return err
	}

	e.AccountID = value
	return nil
}
