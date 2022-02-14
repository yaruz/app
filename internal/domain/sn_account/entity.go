package sn_account

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
)

const (
	EntityType               = "SnAccount"
	PropertySysnameEmail     = "SnAccount.Email"
	PropertySysnameAccountID = "SnAccount.AccountID"
)

// SnAccount is the SnAccount entity
type SnAccount struct {
	*entity.Entity
	ID        uint
	AccountID string
	Email     string
}

var _ entity.Searchable = (*SnAccount)(nil)

func (e SnAccount) EntityType() string {
	return EntityType
}

func (e SnAccount) Validate() error {
	return nil
}

func (e *SnAccount) GetMapNameSysname() map[string]string {
	return map[string]string{
		"Email":     PropertySysnameEmail,
		"AccountID": PropertySysnameAccountID,
	}
}

func (e *SnAccount) SetEmail(ctx context.Context, value string) error {
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

func (e *SnAccount) SetAccountID(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnameAccountID, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.AccountID = value
	return nil
}
