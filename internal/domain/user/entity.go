package user

import (
	"context"
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
)

const (
	EntityType          = "user"
	PropertySysnameName = "name"
)

// User is the user entity
type User struct {
	*entity.Entity
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (e User) EntityType() string {
	return EntityType
}

func (e User) Validate() error {
	return nil
}

func (e *User) SetName(ctx context.Context, value string, langID uint) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnameName, langID)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, langID); err != nil {
		return err
	}

	e.Name = value
	return nil
}
