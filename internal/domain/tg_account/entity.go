package tg_account

import (
	"context"
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
)

const (
	EntityType               = "TgAccount"
	PropertySysnameTgID      = "TgAccount.TgID"
	PropertySysnameCreatedAt = "TgAccount.CreatedAt"
)

var validPropertySysnames = []string{
	PropertySysnameTgID,
	PropertySysnameCreatedAt,
}

// TgAccount is the TgAccount entity
type TgAccount struct {
	*entity.Entity
	ID        uint
	TgID      string
	CreatedAt time.Time `json:"created"`
}

var _ entity.Searchable = (*TgAccount)(nil)

func (e TgAccount) EntityType() string {
	return EntityType
}

func (e TgAccount) Validate() error {
	return nil
}

func (e *TgAccount) GetValidPropertySysnames() []string {
	return validPropertySysnames
}

func (e *TgAccount) GetMapNameSysname() map[string]string {
	return map[string]string{
		"TgID":      PropertySysnameTgID,
		"CreatedAt": PropertySysnameCreatedAt,
	}
}

func (e *TgAccount) SetTgID(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnameTgID, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.TgID = value
	return nil
}

func (e *TgAccount) SetCreatedAt(ctx context.Context, value time.Time) error {
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
