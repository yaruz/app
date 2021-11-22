package user

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
)

const (
	EntityType            = "user"
	PropertySysnameName   = "user.name"
	PropertySysnameAge    = "user.age"
	PropertySysnameHeight = "user.height"
	PropertySysnameWeight = "user.weight"
)

// User is the user entity
type User struct {
	*entity.Entity
	Name   string
	Age    uint
	Height float64
	Weight float64
	//CreatedAt time.Time
	//UpdatedAt time.Time
	//DeletedAt *time.Time
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
		"Name":   PropertySysnameName,
		"Age":    PropertySysnameAge,
		"Height": PropertySysnameHeight,
		"Weight": PropertySysnameWeight,
	}
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

func (e *User) SetAge(ctx context.Context, value uint) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnameAge, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.Age = value
	return nil
}

func (e *User) SetHeight(ctx context.Context, value float64) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnameHeight, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.Height = value
	return nil
}

func (e *User) SetWeight(ctx context.Context, value float64) error {
	prop, err := e.PropertyFinder.GetBySysname(ctx, PropertySysnameWeight, 0)
	if err != nil {
		return err
	}

	if err = e.Entity.SetValueForProperty(prop, value, 0); err != nil {
		return err
	}

	e.Weight = value
	return nil
}
