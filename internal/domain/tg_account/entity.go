package tg_account

import (
	"context"
	"encoding/json"
	"time"

	mtproto_session "github.com/Kalinin-Andrey/mtproto/session"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
)

const (
	EntityType                 = "TgAccount"
	PropertySysnameTgID        = "TgAccount.TgID"
	PropertySysnameAuthSession = "TgAccount.AuthSession"
	PropertySysnameCreatedAt   = "TgAccount.CreatedAt"
)

var validPropertySysnames = []string{
	PropertySysnameTgID,
	PropertySysnameAuthSession,
	PropertySysnameCreatedAt,
}

// TgAccount is the TgAccount entity
type TgAccount struct {
	*entity.Entity
	ID          uint
	TgID        string
	AuthSession *mtproto_session.Session
	CreatedAt   time.Time `json:"created"`
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
		"TgID":        PropertySysnameTgID,
		"AuthSession": PropertySysnameAuthSession,
		"CreatedAt":   PropertySysnameCreatedAt,
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

func (e *TgAccount) SetAuthSession(ctx context.Context, authSession *mtproto_session.Session) error {
	valueb, err := json.Marshal(*authSession)
	if err != nil {
		return err
	}

	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnameAuthSession, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, string(valueb), 0); err != nil {
		return err
	}

	e.AuthSession = authSession
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
