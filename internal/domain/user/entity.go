package user

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"time"
)

const (
	EntityType               = "User"
	PropertySysnameEmail     = "User.Email"
	PropertySysnamePhone     = "User.Phone"
	PropertySysnameFirstName = "User.FirstName"
	PropertySysnameLastName  = "User.LastName"
	PropertySysnameUserName  = "User.UserName"
	PropertySysnameCreatedAt = "User.CreatedAt"
	RelationSysnameTgAccount = "User.TgAccount"
)

var validPropertySysnames = []string{
	PropertySysnameEmail,
	PropertySysnamePhone,
	PropertySysnameFirstName,
	PropertySysnameLastName,
	PropertySysnameUserName,
	PropertySysnameCreatedAt,
}

// User is the user entity
type User struct {
	*entity.Entity
	ID        uint
	Email     string
	Phone     string
	FirstName string
	LastName  string
	UserName  string
	CreatedAt time.Time
}

var _ entity.Searchable = (*User)(nil)

func (e *User) EntityType() string {
	return EntityType
}

func (e *User) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.Email, is.Email),
		validation.Field(&e.Phone, is.E164, validation.Required),
	)
}

func (e *User) GetValidPropertySysnames() []string {
	return validPropertySysnames
}

func (e *User) GetMapNameSysname() map[string]string {
	return map[string]string{
		"Email":     PropertySysnameEmail,
		"Phone":     PropertySysnamePhone,
		"FirstName": PropertySysnameFirstName,
		"LastName":  PropertySysnameLastName,
		"UserName":  PropertySysnameUserName,
		"CreatedAt": PropertySysnameCreatedAt,
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

func (e *User) SetFirstName(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnameFirstName, 1)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 1); err != nil {
		return err
	}

	e.Email = value
	return nil
}

func (e *User) SetLastName(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnameLastName, 1)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 1); err != nil {
		return err
	}

	e.Email = value
	return nil
}

func (e *User) SetUserName(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnameUserName, 1)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 1); err != nil {
		return err
	}

	e.Email = value
	return nil
}

func (e *User) SetPhone(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnamePhone, 1)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 1); err != nil {
		return err
	}

	e.Phone = value
	return nil
}

func (e *User) SetCreatedAt(ctx context.Context, value time.Time) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnameCreatedAt, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.CreatedAt = value
	return nil
}
