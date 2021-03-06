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
	PropertySysnameAccountID = "User.AccountID"
	PropertySysnameCreatedAt = "User.CreatedAt"
	RelationSysnameTgAccount = "User.TgAccount"
)

var validPropertySysnames = []string{
	PropertySysnameEmail,
	PropertySysnamePhone,
	PropertySysnameAccountID,
	PropertySysnameCreatedAt,
}

// User is the user entity
type User struct {
	*entity.Entity
	ID        uint
	AccountID string
	Email     string
	Phone     string
	CreatedAt time.Time `json:"created"`
}

var _ entity.Searchable = (*User)(nil)

func (e *User) EntityType() string {
	return EntityType
}

func (e *User) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.AccountID, validation.Required),
		validation.Field(&e.Email, is.Email, validation.Required),
		validation.Field(&e.Phone, is.E164),
	)
}

func (e *User) GetValidPropertySysnames() []string {
	return validPropertySysnames
}

func (e *User) GetMapNameSysname() map[string]string {
	return map[string]string{
		"Email":     PropertySysnameEmail,
		"Phone":     PropertySysnamePhone,
		"AccountID": PropertySysnameAccountID,
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
