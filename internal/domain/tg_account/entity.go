package tg_account

import (
	"context"
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
)

const (
	EntityType                 = "TgAccount"
	PropertySysname_UserID     = "TgAccount.UserID"
	PropertySysname_AccessHash = "TgAccount.AccessHash"
	PropertySysname_FirstName  = "TgAccount.FirstName"
	PropertySysname_LastName   = "TgAccount.LastName"
	PropertySysname_UserName   = "TgAccount.UserName"
	PropertySysname_Phone      = "TgAccount.Phone"
	PropertySysname_Photo      = "TgAccount.Photo"
	PropertySysname_LangCode   = "TgAccount.LangCode"
	PropertySysname_CreatedAt  = "TgAccount.CreatedAt"
)

var validPropertySysnames = []string{
	PropertySysname_UserID,
	PropertySysname_AccessHash,
	PropertySysname_FirstName,
	PropertySysname_LastName,
	PropertySysname_UserName,
	PropertySysname_Phone,
	PropertySysname_Photo,
	PropertySysname_LangCode,
	PropertySysname_CreatedAt,
}

// TgAccount is the TgAccount entity
type TgAccount struct {
	*entity.Entity
	ID         uint
	UserID     int // в gotd/td это поле int64, но ярус поддерживает только int
	AccessHash int // в gotd/td это поле int64, но ярус поддерживает только int
	FirstName  string
	LastName   string
	UserName   string
	Phone      string
	Photo      string
	LangCode   string
	CreatedAt  time.Time
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
		"UserID":     PropertySysname_UserID,
		"AccessHash": PropertySysname_AccessHash,
		"FirstName":  PropertySysname_FirstName,
		"LastName":   PropertySysname_LastName,
		"UserName":   PropertySysname_UserName,
		"Phone":      PropertySysname_Phone,
		"Photo":      PropertySysname_Photo,
		"LangCode":   PropertySysname_LangCode,
		"CreatedAt":  PropertySysname_CreatedAt,
	}
}

func (e *TgAccount) SetUserID(ctx context.Context, value int) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysname_UserID, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.UserID = value
	return nil
}

func (e *TgAccount) SetAccessHash(ctx context.Context, value int) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysname_AccessHash, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.AccessHash = value
	return nil
}

func (e *TgAccount) SetFirstName(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysname_FirstName, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.FirstName = value
	return nil
}

func (e *TgAccount) SetLastName(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysname_LastName, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.LastName = value
	return nil
}

func (e *TgAccount) SetUserName(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysname_UserName, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.UserName = value
	return nil
}

func (e *TgAccount) SetPhone(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysname_Phone, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.Phone = value
	return nil
}

func (e *TgAccount) SetPhoto(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysname_Photo, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.Photo = value
	return nil
}

func (e *TgAccount) SetLangCode(ctx context.Context, value string) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysname_LangCode, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.LangCode = value
	return nil
}

func (e *TgAccount) SetCreatedAt(ctx context.Context, value time.Time) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysname_CreatedAt, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.CreatedAt = value
	return nil
}
