package advertiser

import (
	"context"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
)

const (
	EntityType          = "Advertiser"
	PropertySysnameName = "Advertiser.Name"
)

var validPropertySysnames = []string{
	PropertySysnameName,
}

// Advertiser is the user entity
type Advertiser struct {
	*entity.Entity
	ID   uint
	Name string
}

var _ entity.Searchable = (*Advertiser)(nil)

func (e Advertiser) EntityType() string {
	return EntityType
}

func (e Advertiser) Validate() error {
	return nil
}

func (e *Advertiser) GetValidPropertySysnames() []string {
	return validPropertySysnames
}

func (e *Advertiser) GetMapNameSysname() map[string]string {
	return map[string]string{
		"Name": PropertySysnameName,
	}
}

func (e *Advertiser) SetName(ctx context.Context, value string, langID uint) error {
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
