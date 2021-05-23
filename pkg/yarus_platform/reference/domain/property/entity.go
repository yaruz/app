package property

import (
	"encoding/json"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

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
	NameSourceID        *uint                    `gorm:"type:bigint null REFERENCES \"text_source\"(id);index" json:"-"`
	DescriptionSourceID *uint                    `gorm:"type:bigint null REFERENCES \"text_source\"(id);index" json:"-"`
	Name                *string                  `gorm:"-" json:"name"`
	Description         *string                  `gorm:"-" json:"description"`
	PropertyTypeID      uint                     `gorm:"type:bigint not null REFERENCES \"property_type\"(id);index" json:"propertyTypeId"`
	PropertyUnitID      uint                     `gorm:"type:bigint not null REFERENCES \"property_unit\"(id);index" json:"propertyUnitId"`
	PropertyViewTypeID  *uint                    `gorm:"type:bigint null REFERENCES \"property_view_type\"(id);index" json:"propertyViewTypeId"`
	PropertyGroupID     *uint                    `gorm:"type:bigint null REFERENCES \"property_group\"(id);index" json:"propertyGroupId"`
	IsSpecific          bool                     `gorm:"type:boolean not null default false;" json:"isSpecific"`
	IsRange             bool                     `gorm:"type:boolean not null default false;" json:"isRange"`
	IsMultiple          bool                     `gorm:"type:boolean not null default false;" json:"isMultiple"`
	SortOrder           uint                     `gorm:"type:smallint;not null;default:9999" json:"sortOrder"`
	OptionsB            datatypes.JSON           `json:"-"`
	Options             []map[string]interface{} `gorm:"-" json:"options"`
	//PropertyUnit        *property_unit.PropertyUnit `gorm:""`

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

func (e Property) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Sysname, validation.Required, validation.Length(2, 100), validation.Match(regexp.MustCompile("^[a-z0-9_]+$"))),
		validation.Field(&e.IsRange, validation.When(e.IsMultiple, validation.Empty)),
		validation.Field(&e.IsMultiple, validation.When(e.IsRange, validation.Empty)),
	)
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
