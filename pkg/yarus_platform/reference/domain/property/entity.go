package property

import (
	"encoding/json"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_group"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_unit"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_view_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/sysname"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	EntityName          = "property"
	TableName           = "property"
	ParseFormateDate    = "2006-01-02"
	ParseFormateDateRFC = time.RFC3339
	ParseFormateTime    = time.RFC3339
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

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
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
		validation.Field(&e.Sysname, sysname.ValidationRules...),
		validation.Field(&e.PropertyTypeID, validation.In(property_type.ValidTypesIDs...)),
		validation.Field(&e.IsRange, validation.When(e.IsMultiple, validation.Empty)),
		validation.Field(&e.IsMultiple, validation.When(e.IsRange, validation.Empty)),
		//validation.Field(&e.PropertyUnitID, validation.When(e.PropertyTypeID != property_type.IDRelation, validation.Required)),
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
				if _, err := GetValueBool(itemVal); err != nil {
					err = errors.Errorf("value type must be a boolean, value: %#v", itemVal)
					break LOOP
				}

			case property_type.IDInt:
				if _, err := GetValueInt(itemVal); err != nil {
					err = errors.Errorf("value type must be an int, value: %#v", itemVal)
					break LOOP
				}

			case property_type.IDFloat:
				if _, err := GetValueFloat(itemVal); err != nil {
					if _, ok := itemVal.(int); !ok {
						err = errors.Errorf("value type must be a float, value: %#v", itemVal)
						break LOOP
					}
				}

			case property_type.IDDate:
				if _, err := GetValueDate(itemVal); err != nil {
					if _, ok := itemVal.(int); !ok {
						err = errors.Errorf("value type must be a date, value: %#v", itemVal)
						break LOOP
					}
				}

			case property_type.IDTime:
				if _, err := GetValueTime(itemVal); err != nil {
					err = errors.Errorf("value type must be a timestamp, value: %#v", itemVal)
				}

			case property_type.IDText:
				if _, err := GetValueText(itemVal); err != nil {
					err = errors.Errorf("value type must be a text, value: %#v", itemVal)
					break LOOP
				}
			}
		}
	}
	return err
}

func GetValueBool(value interface{}) (bool, error) {
	res, ok := value.(bool)
	if !ok {
		return res, errors.Errorf("can not cast to a boolean value, value: %#v", value)
	}
	return res, nil
}

func GetValueInt(value interface{}) (int, error) {
	var res int

	valInt, okInt := value.(int)
	valUint, okUint := value.(uint)
	valFloat, okFloat := value.(float64) // после анмаршаллинга из JSON тип float64

	if okInt {
		res = valInt
	} else if okUint {
		if uint(int(valUint)) != valUint {
			return res, errors.Errorf("Can not cast value of PropertyValue to int from uint. Value = %v.", value)
		}
		res = int(valUint)
	} else if okFloat {
		if float64(int(valFloat)) != valFloat {
			return res, errors.Errorf("Can not cast value of PropertyValue to int from float64. Value = %v.", value)
		}
		res = int(valFloat)
	} else {
		return res, errors.Errorf("Can not cast value of PropertyValue to int. Value = %v.", value)
	}
	return res, nil
}

func GetRelationItemValue(value interface{}) (uint, error) {
	var res uint

	valInt, okInt := value.(int)
	valUint, okUint := value.(uint)
	valFloat, okFloat := value.(float64) // на всякий..

	switch {
	case okUint:
		res = valUint
	case okInt:
		if int(uint(valInt)) != valInt {
			return res, errors.Errorf("Can not cast value of Relation to uint from int. Value = %v.", value)
		}

		res = uint(valInt)
	case okFloat:
		if float64(uint(valFloat)) != valFloat {
			return res, errors.Errorf("Can not cast value of Relation to uint from float64. Value = %v.", value)
		}

		res = uint(valInt)
	default:
		return res, errors.Errorf("Can not cast value of Relation to uint. Value = %v.", value)
	}

	return res, nil
}

func GetRelationValue(value interface{}) ([]uint, error) {
	var res []uint

	valInt, okInt := value.([]int)
	valUint, okUint := value.([]uint)
	valFloat, okFloat := value.([]float64)             // на всякий..
	valInterface, okInterface := value.([]interface{}) // после анмаршаллинга из JSON тип []interface{}

	switch {
	case okUint:
		res = valUint
	case okInt:
		res = make([]uint, 0, len(valInt))

		for _, i := range valInt {

			if int(uint(i)) != i {
				return res, errors.Errorf("Can not cast value of Relation to []uint from []int. Value = %v.", value)
			}
			res = append(res, uint(i))
		}
	case okFloat:
		res = make([]uint, 0, len(valInt))

		for _, i := range valFloat {

			if float64(uint(i)) != i {
				return res, errors.Errorf("Can not cast value of Relation to []uint from []float64. Value = %v.", value)
			}
			res = append(res, uint(i))
		}
	case okInterface:
		res = make([]uint, 0, len(valInt))

		for _, val := range valInterface {
			valInt, err := GetValueInt(val)
			if err != nil {
				return nil, err
			}
			res = append(res, uint(valInt))
		}
	default:
		return res, errors.Errorf("Can not cast value of Relation to []uint. Value = %v.", value)
	}

	return res, nil
}

func GetValueFloat(value interface{}) (float64, error) {
	res, ok := value.(float64)
	if !ok {

		resInt, ok := value.(int)
		if !ok {
			return res, errors.Errorf("Can not cast value of PropertyValue to float. Value = %v.", value)
		}
		res = float64(resInt)
	}
	return res, nil
}

func GetValueDate(value interface{}) (time.Time, error) {
	var res time.Time
	var t time.Time
	var err error

	t, okTime := value.(time.Time)
	vStr, okStr := value.(string)
	vInt64, okInt64 := value.(int64)
	vInt, okInt := value.(int)

	if okTime {
	} else if okStr {

		if t, err = time.Parse(ParseFormateDateRFC, vStr); err != nil {
			if t, err = time.Parse(ParseFormateDate, vStr); err != nil {
				return res, errors.Wrapf(err, "Can not pars string value by ParseFormateDate to date. Value = %v.", value)
			}
		}
	} else if okInt || okInt64 {

		if okInt {
			vInt64 = int64(vInt)
		}

		t = time.Unix(vInt64, 0)
	} else {
		return res, errors.Errorf("Can not cast value of PropertyValue to string (as a string format of a date) or int (as a timestamp). Value = %v.", value)
	}

	y, m, d := t.Date()
	res = time.Date(y, m, d, 0, 0, 0, 0, t.Location())

	return res, nil
}

func GetValueTime(value interface{}) (time.Time, error) {
	var t time.Time
	var err error

	t, okTime := value.(time.Time)
	vStr, okStr := value.(string)
	vInt64, okInt64 := value.(int64)
	vInt, okInt := value.(int)

	if okTime {
	} else if okStr {
		t, err = time.Parse(ParseFormateTime, vStr)
		if err != nil {
			return t, errors.Wrapf(err, "Can not pars string value by ParseFormateDate to date. Value = %v.", value)
		}
	} else if okInt || okInt64 {

		if okInt {
			vInt64 = int64(vInt)
		}

		t = time.Unix(vInt64, 0)
	} else {
		return t, errors.Errorf("Can not cast value of PropertyValue to string (as a string format of a date) or int (as a timestamp). Value = %v.", value)
	}

	return t, nil
}

func GetValueText(value interface{}) (string, error) {
	res, ok := value.(string)
	if !ok {
		return res, errors.Errorf("can not cast to a string value, value: %#v", value)
	}
	return res, nil
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
