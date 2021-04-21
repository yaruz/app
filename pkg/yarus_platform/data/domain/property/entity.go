package property

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

const (
	EntityName = "property"
	TableName  = "property"
)

// Property ...
type Property struct {
	ID                  uint                     `gorm:"type:bigserial;primaryKey" json:"id"`
	Sysname             string                   `gorm:"type:varchar(100) not null;unique;index" json:"sysname"`
	NameSourceID        uint                     `sql:"type:bigint not null REFERENCES \"text_source\"(id)" gorm:"index" json:"-"`
	DescriptionSourceID uint                     `sql:"type:bigint not null REFERENCES \"text_source\"(id)" gorm:"index" json:"-"`
	Name                string                   `gorm:"-" json:"name"`
	Description         string                   `gorm:"-" json:"description"`
	PropertyTypeID      uint                     `sql:"type:bigint not null REFERENCES \"property_type\"(id)" gorm:"index" json:"propertyTypeId"`
	PropertyUnitID      uint                     `sql:"type:bigint not null REFERENCES \"property_unit\"(id)" gorm:"index" json:"propertyUnitId"`
	PropertyViewTypeID  uint                     `sql:"type:bigint REFERENCES \"property_view_type\"(id)" gorm:"index" json:"propertyViewTypeId"`
	PropertyGroupID     uint                     `sql:"type:bigint REFERENCES \"property_group\"(id)" gorm:"index" json:"propertyGroupId"`
	IsSpecific          bool                     `gorm:"type:boolean not null default false;" json:"isSpecific"`
	IsRange             bool                     `gorm:"type:boolean not null default false;" json:"isRange"`
	IsMultiple          bool                     `gorm:"type:boolean not null default false;" json:"isMultiple"`
	SortOrder           uint                     `gorm:"type:smallint not null default 9999" json:"sortOrder"`
	OptionsB            datatypes.JSON           `json:"-"`
	Options             []map[string]interface{} `gorm:"-" json:"options"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (e *Property) TableName() string {
	return TableName
}

// New func is a constructor for the Property
func New() *Property {
	return &Property{}
}

func (e *Property) AfterFind() error {
	return e.optionsB2Options()
}

func (e *Property) BeforeSave() error {
	return e.options2OptionsB()
}

func (e *Property) optionsB2Options() error {
	jsonb, err := e.OptionsB.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonb, &e.Options)
}

func (e *Property) options2OptionsB() error {
	jsonb, err := json.Marshal(&e.Options)
	if err != nil {
		return err
	}
	return e.OptionsB.UnmarshalJSON(jsonb)
}
