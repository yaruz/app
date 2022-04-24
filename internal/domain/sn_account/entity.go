package sn_account

import (
	"context"
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
)

const (
	EntityType          = "SNAccount"
	PropertySysnameSNID = "SNAccount.SNID"
)

var validPropertySysnames = []string{
	PropertySysnameSNID,
}

// SNAccount is the SNAccount entity
type SNAccount struct {
	*entity.Entity
	ID        uint
	TypeID    uint
	SNID      string
	CreatedAt time.Time `json:"created"`
}

var _ entity.Searchable = (*SNAccount)(nil)

func (e SNAccount) EntityType() string {
	return EntityType
}

func (e SNAccount) Validate() error {
	return nil
}

func (e *SNAccount) GetValidPropertySysnames() []string {
	return validPropertySysnames
}

func (e *SNAccount) GetMapNameSysname() map[string]string {
	return map[string]string{
		"SNID": PropertySysnameSNID,
	}
}

func (e *SNAccount) SetAccountID(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnameSNID, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.SNID = value
	return nil
}
