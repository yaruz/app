package property

import (
	"github.com/yaruz/app/pkg/yarus_platform/config"
)

type Configs map[string]Config

type Config struct {
	PropertyTypeID     uint
	PropertyUnitID     *uint
	PropertyViewTypeID *uint
	PropertyGroupID    *uint
	IsSpecific         bool
	IsRange            bool
	IsMultiple         bool
	SortOrder          uint
	Options            []map[string]interface{}
	Texts              map[string]config.NameAndDescriptionText
}
