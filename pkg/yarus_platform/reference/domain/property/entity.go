package property

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_group"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_unit"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_view_type"

	"github.com/pkg/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"

	"gorm.io/datatypes"
)

const (
	EntityName       = "property"
	TableName        = "property"
	ParseFormateDate = "2006-01-02"
)

// Property ...
type Property struct {
	ID                  uint                                 `gorm:"type:bigserial;primaryKey" json:"id"`
	Sysname             string                               `gorm:"type:varchar(100) not null;unique;index" json:"sysname"`
	NameSourceID        *uint                                `gorm:"type:bigint null REFERENCES \"text_source\"(id);index" json:"-"`
	DescriptionSourceID *uint                                `gorm:"type:bigint null REFERENCES \"text_source\"(id);index" json:"-"`
	Name                *string                              `gorm:"-" json:"name"`
	Description         *string                              `gorm:"-" json:"description"`
	PropertyTypeID      uint                                 `gorm:"type:bigint not null REFERENCES \"property_type\"(id);index" json:"propertyTypeId"`
	PropertyUnitID      *uint                                `gorm:"type:bigint null REFERENCES \"property_unit\"(id);index" json:"propertyUnitId"`
	PropertyViewTypeID  *uint                                `gorm:"type:bigint null REFERENCES \"property_view_type\"(id);index" json:"propertyViewTypeId"`
	PropertyGroupID     *uint                                `gorm:"type:bigint null REFERENCES \"property_group\"(id);index" json:"propertyGroupId"`
	IsSpecific          bool                                 `gorm:"type:boolean not null default false;" json:"isSpecific"`
	IsRange             bool                                 `gorm:"type:boolean not null default false;" json:"isRange"`
	IsMultiple          bool                                 `gorm:"type:boolean not null default false;" json:"isMultiple"`
	SortOrder           uint                                 `gorm:"type:smallint;not null;default:9999" json:"sortOrder"`
	OptionsB            datatypes.JSON                       `json:"-"`
	Options             []map[string]interface{}             `gorm:"-" json:"options"`
	PropertyType        *property_type.PropertyType          `json:"propertyType"`
	PropertyViewType    *property_view_type.PropertyViewType `json:"propertyViewType"`
	PropertyUnit        *property_unit.PropertyUnit          `json:"propertyUnit"`
	PropertyGroup       *property_group.PropertyGroup        `json:"propertyGroup"`

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
		validation.Field(&e.Sysname, domain.SysnameValidationRules...),
		validation.Field(&e.PropertyTypeID, validation.In(property_type.ValidTypesIDs...)),
		validation.Field(&e.IsRange, validation.When(e.IsMultiple, validation.Empty)),
		validation.Field(&e.IsMultiple, validation.When(e.IsRange, validation.Empty)),
		validation.Field(&e.PropertyUnitID, validation.When(e.PropertyTypeID != property_type.IDRelation, validation.Required)),
		validation.Field(&e.Options, validation.By(e.optionsValidate)),
	)
}

func (e Property) optionsValidate(value interface{}) (err error) {
	v, ok := value.([]map[string]interface{})
	if !ok {
		return errors.New("Property.Options must be []map[string]interface{} type")
	}

LOOP:
	for _, item := range v {

		for _, itemVal := range item {

			switch e.PropertyTypeID {

			case property_type.IDBoolean:
				if _, ok := itemVal.(bool); !ok {
					err = errors.Errorf("value type must be a boolean, value: %#v", itemVal)
					break LOOP
				}

			case property_type.IDInt:
				if _, ok := itemVal.(int); !ok {
					err = errors.Errorf("value type must be an int, value: %#v", itemVal)
					break LOOP
				}

			case property_type.IDBigint, property_type.IDRelation:
				if _, ok := itemVal.(int); !ok {
					if _, ok := itemVal.(int64); !ok {
						err = errors.Errorf("value type must be a bigint, value: %#v", itemVal)
						break LOOP
					}
				}

			case property_type.IDFloat:
				if _, ok := itemVal.(float64); !ok {
					if _, ok := itemVal.(int); !ok {
						err = errors.Errorf("value type must be a float, value: %#v", itemVal)
						break LOOP
					}
				}

			case property_type.IDDate:
				s, ok := itemVal.(string)
				if !ok {
					err = errors.Errorf("value type must be a date, value: %#v", itemVal)
					break LOOP
				}
				//d := &datatypes.Date{}
				var t time.Time
				t, err = time.Parse(ParseFormateDate, s)
				if err != nil {
					err = errors.Errorf("value type must be a date, value: %#v", itemVal)
					break LOOP
				}
				d := datatypes.Date(t)
				dv, err := d.Value()
				if err != nil {
					err = errors.Errorf("value type must be a date, value: %#v", itemVal)
					break LOOP
				}
				fmt.Printf("%#v", dv)

			case property_type.IDTimestamp:
				if _, ok := itemVal.(int); !ok {
					if _, ok := itemVal.(int64); !ok {
						err = errors.Errorf("value type must be a timestamp, value: %#v", itemVal)
						break LOOP
					}
				}

			case property_type.IDText:
				if _, ok := itemVal.(string); !ok {
					err = errors.Errorf("value type must be a text, value: %#v", itemVal)
					break LOOP
				}
			}
		}
	}
	return err
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
