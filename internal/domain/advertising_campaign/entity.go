package advertising_campaign

import (
	"context"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
)

const (
	EntityType          = "AdvertisingCampaign"
	PropertySysnameName = "AdvertisingCampaign.Name"
)

var validPropertySysnames = []string{
	PropertySysnameName,
}

// AdvertisingCampaign is the user entity
type AdvertisingCampaign struct {
	*entity.Entity
	ID   uint
	Name string
}

var _ entity.Searchable = (*AdvertisingCampaign)(nil)

func (e AdvertisingCampaign) EntityType() string {
	return EntityType
}

func (e AdvertisingCampaign) Validate() error {
	return nil
}

func (e *AdvertisingCampaign) GetValidPropertySysnames() []string {
	return validPropertySysnames
}

func (e *AdvertisingCampaign) GetMapNameSysname() map[string]string {
	return map[string]string{
		"Name": PropertySysnameName,
	}
}

func (e *AdvertisingCampaign) SetName(ctx context.Context, value string, langID uint) error {
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
