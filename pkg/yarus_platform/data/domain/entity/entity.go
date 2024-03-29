package entity

import (
	"context"
	"sort"

	"github.com/pkg/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/bool_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/date_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/float_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/int_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/time_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/utext_value"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
)

const (
	EntityName = "entity"
	TableName  = "entity"

	FieldName_ID           = "ID"
	FieldName_EntityTypeID = "EntityTypeID"
	FieldName_EntityType   = "EntityType"
)

type LangFinder interface {
	GetCodes(ctx context.Context) ([]string, error)
	GetCodesEmptyInterfaceSlice(ctx context.Context) ([]interface{}, error)
	GetMapCodeID(ctx context.Context) (map[string]uint, error)
	GetMapIDCode(ctx context.Context) (map[uint]string, error)
	GetMapIDCfgname(ctx context.Context) (map[uint]string, error)
	GetIDByCode(ctx context.Context, code string) (uint, error)
	GetCfgnameByID(ctx context.Context, id uint) (string, error)
}

type EntityTypeFinder interface {
	Get(ctx context.Context, id uint) (*entity_type.EntityType, error)
	GetBySysname(ctx context.Context, sysname string, langID uint) (*entity_type.EntityType, error)
	GetSysnames(ctx context.Context) ([]string, error)
	GetSysnamesEmptyInterfaceSlice(ctx context.Context) ([]interface{}, error)
	GetMapSysnameID(ctx context.Context) (map[string]uint, error)
	GetMapIDSysname(ctx context.Context) (map[uint]string, error)
	GetIDBySysname(ctx context.Context, sysname string) (uint, error)
	GetSysnameByID(ctx context.Context, id uint) (string, error)
}

type PropertyFinder interface {
	Get(ctx context.Context, id uint) (*property.Property, error)
	GetBySysname(ctx context.Context, sysname string, langID uint) (*property.Property, error)
	GetSysnames(ctx context.Context) ([]string, error)
	GetSysnamesEmptyInterfaceSlice(ctx context.Context) ([]interface{}, error)
	GetMapSysnameID(ctx context.Context) (map[string]uint, error)
	GetMapIDSysname(ctx context.Context) (map[uint]string, error)
	GetMapSysnameTypeID(ctx context.Context) (map[string]uint, error)
	GetMapIDTypeID(ctx context.Context) (map[uint]uint, error)
	GetIDBySysname(ctx context.Context, sysname string) (uint, error)
	GetSysnameByID(ctx context.Context, id uint) (string, error)
	GetTypeIDBySysname(ctx context.Context, sysname string) (uint, error)
	GetTypeIDByID(ctx context.Context, id uint) (uint, error)
}

type Searchable interface {
	GetValidPropertySysnames() []string
	GetMapNameSysname() map[string]string
}

// Entity ...
// Значения свойст:
// 	пишем в соотв. слайс из ...Values
// 	храним в соотв. таблице ..._value
// 	PropertiesValues - красивый и удобный список со свойствами, только для чтения, инициализируется в сервисе entity.Service.EntityInit()
// 	RelationsValues - красивый и удобный список со связями - // - // -
//
type Entity struct {
	ID               uint                     `gorm:"type:bigint;primaryKey" json:"id"`
	EntityTypeID     uint                     `gorm:"type:bigint not null;index" json:"entityTypeID"`
	PropertiesValues map[uint]PropertyValue   `gorm:"-" json:"propertiesValues"`
	RelationsValues  map[uint]RelationValue   `gorm:"-" json:"relationsValues"`
	BoolValues       []bool_value.BoolValue   `json:"-"`
	IntValues        []int_value.IntValue     `json:"-"`
	FloatValues      []float_value.FloatValue `json:"-"`
	DateValues       []date_value.DateValue   `json:"-"`
	TimeValues       []time_value.TimeValue   `json:"-"`
	TextValues       []text_value.TextValue   `json:"-"`
	UTextValues      []utext_value.UTextValue `json:"-"`
	PropertyFinder   PropertyFinder           `gorm:"-" json:"-"`
}

func (e *Entity) TableName() string {
	return TableName
}

// New func is a constructor for the Entity
func New() *Entity {
	return &Entity{}
}

func (e Entity) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.EntityTypeID, validation.Required),
	)
}

// Присваивает PropertyValue
func (e *Entity) setPropertyValue(propertyValue *PropertyValue) {
	propertyID := propertyValue.Property.ID

	if e.PropertiesValues == nil {
		e.PropertiesValues = make(map[uint]PropertyValue, 1)
	}

	e.PropertiesValues[propertyID] = *propertyValue
}

// Присваивает RelationValue
func (e *Entity) setRelationValue(relationValue *RelationValue) {
	propertyID := relationValue.Property.ID

	if e.RelationsValues == nil {
		e.RelationsValues = make(map[uint]RelationValue, 1)
	}

	e.RelationsValues[propertyID] = *relationValue
}

// Удаляет все как значения свойств, так и значения связей для заданного propertyID
func (e *Entity) DeletePropertyValues(propertyID uint) {

	propertyValue, propOk := e.PropertiesValues[propertyID]
	_, relOk := e.RelationsValues[propertyID]

	if propOk {
		delete(e.PropertiesValues, propertyID)

		switch propertyValue.Property.PropertyTypeID {
		case property_type.IDBoolean:
			e.deletePropertyBoolValues(propertyID)
		case property_type.IDInt:
			e.deletePropertyIntValues(propertyID)
		case property_type.IDFloat:
			e.deletePropertyFloatValues(propertyID)
		case property_type.IDDate:
			e.deletePropertyDateValues(propertyID)
		case property_type.IDTime:
			e.deletePropertyTimeValues(propertyID)
		case property_type.IDText:
			e.deletePropertyTextValues(propertyID)
		case property_type.IDUText:
			e.deletePropertyUTextValues(propertyID)
		}
	} else if relOk {
		delete(e.RelationsValues, propertyID)
		e.deletePropertyIntValues(propertyID)
	}
}

// Удаляет все значения BoolValues для заданного propertyID
func (e *Entity) deletePropertyBoolValues(propertyID uint) {
	items := make([]bool_value.BoolValue, 0, len(e.BoolValues))
	for i := range e.BoolValues {
		if e.BoolValues[i].PropertyID != propertyID {
			items = append(items, e.BoolValues[i])
		}
	}
	e.BoolValues = items
}

// Удаляет все значения IntValues для заданного propertyID
func (e *Entity) deletePropertyIntValues(propertyID uint) {
	items := make([]int_value.IntValue, 0, len(e.IntValues))
	for i := range e.IntValues {
		if e.IntValues[i].PropertyID != propertyID {
			items = append(items, e.IntValues[i])
		}
	}
	e.IntValues = items
}

// Удаляет все значения FloatValues для заданного propertyID
func (e *Entity) deletePropertyFloatValues(propertyID uint) {
	items := make([]float_value.FloatValue, 0, len(e.FloatValues))
	for i := range e.FloatValues {
		if e.FloatValues[i].PropertyID != propertyID {
			items = append(items, e.FloatValues[i])
		}
	}
	e.FloatValues = items
}

// Удаляет все значения DateValues для заданного propertyID
func (e *Entity) deletePropertyDateValues(propertyID uint) {
	items := make([]date_value.DateValue, 0, len(e.DateValues))
	for i := range e.DateValues {
		if e.DateValues[i].PropertyID != propertyID {
			items = append(items, e.DateValues[i])
		}
	}
	e.DateValues = items
}

// Удаляет все значения TimeValues для заданного propertyID
func (e *Entity) deletePropertyTimeValues(propertyID uint) {
	items := make([]time_value.TimeValue, 0, len(e.TimeValues))
	for i := range e.TimeValues {
		if e.TimeValues[i].PropertyID != propertyID {
			items = append(items, e.TimeValues[i])
		}
	}
	e.TimeValues = items
}

// Удаляет все значения TextValues для заданного propertyID
func (e *Entity) deletePropertyTextValues(propertyID uint) {
	items := make([]text_value.TextValue, 0, len(e.TextValues))
	for i := range e.TextValues {
		if e.TextValues[i].PropertyID != propertyID {
			items = append(items, e.TextValues[i])
		}
	}
	e.TextValues = items
}

// Удаляет все значения UTextValues для заданного propertyID
func (e *Entity) deletePropertyUTextValues(propertyID uint) {
	items := make([]utext_value.UTextValue, 0, len(e.UTextValues))
	for i := range e.UTextValues {
		if e.UTextValues[i].PropertyID != propertyID {
			items = append(items, e.UTextValues[i])
		}
	}
	e.UTextValues = items
}

// Возвращает количество всех значений всех свойств
func (e *Entity) PropertiesValuesCount() int {
	return len(e.BoolValues) + len(e.IntValues) + len(e.FloatValues) + len(e.DateValues) + len(e.TimeValues) + len(e.TextValues) + len(e.UTextValues)
}

// Возвращает слайс ID всех свойств, значения которых заданны для entity
func (e *Entity) GetPropertiesIDs() []interface{} {
	propertiesIdsMap := make(map[uint]struct{}, e.PropertiesValuesCount())

	for _, val := range e.BoolValues {
		propertiesIdsMap[val.PropertyID] = struct{}{}
	}

	for _, val := range e.IntValues {
		propertiesIdsMap[val.PropertyID] = struct{}{}
	}

	for _, val := range e.FloatValues {
		propertiesIdsMap[val.PropertyID] = struct{}{}
	}

	for _, val := range e.DateValues {
		propertiesIdsMap[val.PropertyID] = struct{}{}
	}

	for _, val := range e.TimeValues {
		propertiesIdsMap[val.PropertyID] = struct{}{}
	}

	for _, val := range e.TextValues {
		propertiesIdsMap[val.PropertyID] = struct{}{}
	}

	for _, val := range e.UTextValues {
		propertiesIdsMap[val.PropertyID] = struct{}{}
	}

	propertiesIds := make([]interface{}, 0, len(propertiesIdsMap))
	for id := range propertiesIdsMap {
		propertiesIds = append(propertiesIds, id)
	}

	return propertiesIds
}

// Возвращает слайс ID значений связей для заданного propertyID
func (e *Entity) GetRelationValues(propertyID uint) []uint {
	values := make([]uint, 0, 1)

	for _, intValue := range e.IntValues {
		if intValue.PropertyID == propertyID {
			values = append(values, uint(intValue.Value))
		}
	}
	return values
}

// Инициализация PropertiesValues и RelationsValues из всех слайсов значений свойств ...Values
// Запускаем после запуска Service.propertiesInit()
func (e *Entity) propertiesValuesInit(langID uint) error {

	for _, propertyValue := range e.BoolValues {
		if err := e.setValueToExistingPropertyValue(propertyValue.PropertyID, propertyValue.Value, langID); err != nil {
			return err
		}
	}

	for _, propertyValue := range e.IntValues {
		if err := e.setValueToExistingPropertyValue(propertyValue.PropertyID, propertyValue.Value, langID); err != nil {
			return err
		}
	}

	for _, propertyValue := range e.FloatValues {
		if err := e.setValueToExistingPropertyValue(propertyValue.PropertyID, propertyValue.Value, langID); err != nil {
			return err
		}
	}

	for _, propertyValue := range e.DateValues {
		if err := e.setValueToExistingPropertyValue(propertyValue.PropertyID, propertyValue.Value, langID); err != nil {
			return err
		}
	}

	for _, propertyValue := range e.TimeValues {
		if err := e.setValueToExistingPropertyValue(propertyValue.PropertyID, propertyValue.Value, langID); err != nil {
			return err
		}
	}

	for _, propertyValue := range e.TextValues {
		if err := e.setValueToExistingPropertyValue(propertyValue.PropertyID, propertyValue.Value, langID); err != nil {
			return err
		}
	}

	for _, propertyValue := range e.UTextValues {
		if err := e.setValueToExistingPropertyValue(propertyValue.PropertyID, propertyValue.Value, langID); err != nil {
			return err
		}
	}

	return nil
}

// Только для метода propertiesValuesInit()
// Уже к существующим после иннициализации экземплярам RelationValue и PropertyValue присваивает значение value
func (e *Entity) setValueToExistingPropertyValue(propertyID uint, value interface{}, langID uint) error {
	propertyValue, propOk := e.PropertiesValues[propertyID]
	relationValue, relOk := e.RelationsValues[propertyID]

	switch {
	case relOk:
		if err := relationValue.AddValueByInterface(value); err != nil {
			return errors.Wrapf(err, "Can not add value to RelationValue. Property ID = %v; Value = %v.", propertyID, value)
		}
		e.RelationsValues[propertyID] = relationValue
	case propOk:
		if err := propertyValue.SetValue(value, langID); err != nil {
			return errors.Wrapf(err, "Can not set value to PropertyValue. Property ID = %v; Value = %v.", propertyID, value)
		}
		e.PropertiesValues[propertyID] = propertyValue
	default:
		return errors.Errorf("Property ID = %v not found.", propertyID)
	}
	return nil
}

// Создаёт и присваивает новые значения PropertyValue
// value - значение, не ссылка
func (e *Entity) SetValueForProperty(property *property.Property, value interface{}, langID uint) error {
	e.DeletePropertyValues(property.ID)
	return e.AddValueForProperty(property, value, langID)
}

// Пока не дописаны множественные значения, работает как Set()
func (e *Entity) AddValueForProperty(prop *property.Property, value interface{}, langID uint) error {
	if value == nil {
		return yaruserror.ErrEmptyParams
	}
	propertyID := prop.ID

	if e.PropertiesValues == nil {
		e.PropertiesValues = make(map[uint]PropertyValue, 1)
	}

	propertyValue, err := newPropertyValue(prop, value, langID)
	if err != nil {
		return err
	}

	// todo пока нет работы с множественными значениями
	if propValue, ok := e.PropertiesValues[propertyID]; ok {
		if propertyValue.LangID != langID {
			errors.Errorf("trying to assign a property value with languageID = %v to value with languageID = %v", langID, propertyValue.LangID)
		}
		propValue.SetValue(value, langID)
		*propertyValue = propValue
	}
	// todo работы с массивами пока нет, соотв. добавления тоже
	e.setPropertyValue(propertyValue)

	switch prop.PropertyTypeID {
	case property_type.IDBoolean:
		value, _ := property.GetValueBool(propertyValue.Value)
		e.BoolValues = append(e.BoolValues, bool_value.BoolValue{
			EntityID:   e.ID,
			PropertyID: propertyID,
			Value:      value,
		})
	case property_type.IDInt:
		value, _ := property.GetValueInt(propertyValue.Value)
		e.IntValues = append(e.IntValues, int_value.IntValue{
			EntityID:   e.ID,
			PropertyID: propertyID,
			Value:      value,
		})
	case property_type.IDFloat:
		value, _ := property.GetValueFloat(propertyValue.Value)
		e.FloatValues = append(e.FloatValues, float_value.FloatValue{
			EntityID:   e.ID,
			PropertyID: propertyID,
			Value:      value,
		})
	case property_type.IDDate:
		value, _ := property.GetValueDate(propertyValue.Value)
		e.DateValues = append(e.DateValues, date_value.DateValue{
			EntityID:   e.ID,
			PropertyID: propertyID,
			Value:      value,
		})
	case property_type.IDTime:
		value, _ := property.GetValueTime(propertyValue.Value)
		e.TimeValues = append(e.TimeValues, time_value.TimeValue{
			EntityID:   e.ID,
			PropertyID: propertyID,
			Value:      value,
		})
	case property_type.IDText:
		value, _ := property.GetValueText(propertyValue.Value)
		e.TextValues = append(e.TextValues, text_value.TextValue{
			EntityID:   e.ID,
			LangID:     propertyValue.LangID,
			PropertyID: propertyID,
			Value:      value,
		})
	case property_type.IDUText:
		value, _ := property.GetValueText(propertyValue.Value)
		e.UTextValues = append(e.UTextValues, utext_value.UTextValue{
			EntityID:   e.ID,
			PropertyID: propertyID,
			Value:      value,
		})
	}

	return nil
}

// Создаёт и присваивает новые значения RelationValue
// value - значение, не ссылка, []uint - IDs связанных сущностей
func (e *Entity) SetValueForRelation(relation *entity_type.Relation, value uint) error {
	e.DeletePropertyValues(relation.Property.ID)
	return e.AddValueForRelation(relation, value)
}

func (e *Entity) SetValuesForRelation(relation *entity_type.Relation, values []uint) error {
	e.DeletePropertyValues(relation.Property.ID)
	return e.AddValuesForRelation(relation, values)
}

func (e *Entity) AddValueForRelation(relation *entity_type.Relation, value uint) error {
	if value == 0 {
		return yaruserror.ErrEmptyParams
	}
	propertyID := relation.Property.ID

	if e.RelationsValues == nil {
		e.RelationsValues = make(map[uint]RelationValue, 1)
	}

	relationValue, err := newRelationValue(relation, []uint{})
	if err != nil {
		return err
	}

	if _, ok := e.RelationsValues[propertyID]; !ok && len(e.GetRelationValues(propertyID)) > 0 {
		return errors.Wrapf(yaruserror.ErrNotSet, "RelationsValues fo relation ID = %v not init", propertyID)
	}

	if relValue, ok := e.RelationsValues[propertyID]; ok {
		*relationValue = relValue
	}

	intValue, err := property.GetValueInt(value)
	if err != nil {
		return err
	}

	if err := relationValue.AddValue(value); err != nil {
		return err
	}

	e.IntValues = append(e.IntValues, int_value.IntValue{
		EntityID:   e.ID,
		PropertyID: propertyID,
		Value:      intValue,
	})

	e.setRelationValue(relationValue)
	return nil
}

func (e *Entity) AddValuesForRelation(relation *entity_type.Relation, values []uint) error {
	if values == nil || len(values) == 0 {
		return yaruserror.ErrEmptyParams
	}
	propertyID := relation.Property.ID
	alreadyExists := make(map[int]interface{})

	if e.RelationsValues == nil {
		e.RelationsValues = make(map[uint]RelationValue, 1)
	}

	relationValue, err := newRelationValue(relation, []uint{})
	if err != nil {
		return err
	}

	if _, ok := e.RelationsValues[propertyID]; !ok && len(e.GetRelationValues(propertyID)) > 0 {
		return errors.Wrapf(yaruserror.ErrNotSet, "RelationsValues fo relation ID = %v not init", propertyID)
	}

	if relValue, ok := e.RelationsValues[propertyID]; ok {
		*relationValue = relValue
	}

	intValues := make([]int, 0, len(values))
	for _, value := range values {
		intValue, err := property.GetValueInt(value)
		if err != nil {
			return err
		}
		intValues = append(intValues, intValue)
	}

	if err = relationValue.AddValues(values); err != nil {
		if !errors.As(err, yaruserror.ErrAlreadyExistsList{}) {
			return err
		}
		errList, ok := err.(yaruserror.ErrAlreadyExistsList)
		if ok {
			alreadyExists = errList.List()
		}
	}

	for i, intValue := range intValues {
		if _, ok := alreadyExists[i]; ok {
			continue
		}

		e.IntValues = append(e.IntValues, int_value.IntValue{
			EntityID:   e.ID,
			PropertyID: propertyID,
			Value:      intValue,
		})
	}

	e.setRelationValue(relationValue)

	if errors.Is(err, yaruserror.ErrAlreadyExists) {
		return err
	}

	return nil
}

func (e *Entity) RemoveValueForRelation(relation *entity_type.Relation, value uint) error {
	if value == 0 {
		return yaruserror.ErrEmptyParams
	}
	propertyID := relation.Property.ID

	if len(e.GetRelationValues(propertyID)) == 0 {
		return errors.Wrapf(yaruserror.ErrNotSet, "RelationsValues fo relation ID = %v is empty", propertyID)
	}

	if e.RelationsValues == nil || len(e.RelationsValues) == 0 {
		return errors.Wrapf(yaruserror.ErrNotSet, "RelationsValues fo relation ID = %v not init", propertyID)
	}

	relationValue, ok := e.RelationsValues[propertyID]
	if !ok {
		return errors.Wrapf(yaruserror.ErrNotFound, "RelationsValues for relation ID = %v not found", propertyID)
	}

	intValue, err := property.GetValueInt(value)
	if err != nil {
		return err
	}

	if err := relationValue.RemoveValue(value); err != nil {
		return err
	}

	if len(relationValue.Value) == 0 {
		e.DeletePropertyValues(propertyID)
	} else {
		e.setRelationValue(&relationValue)

		for i := range e.IntValues {
			if e.IntValues[i].PropertyID == propertyID && e.IntValues[i].Value == intValue {
				e.IntValues = append(e.IntValues[:i], e.IntValues[i+1:]...)
				break
			}
		}
	}

	return nil
}

func (e *Entity) RemoveValuesForRelation(relation *entity_type.Relation, values []uint) error {
	if values == nil || len(values) == 0 {
		return yaruserror.ErrEmptyParams
	}
	propertyID := relation.Property.ID
	notFound := make(map[int]interface{})

	if len(e.GetRelationValues(propertyID)) == 0 {
		return errors.Wrapf(yaruserror.ErrNotSet, "RelationsValues fo relation ID = %v is empty", propertyID)
	}

	if e.RelationsValues == nil || len(e.RelationsValues) == 0 {
		return errors.Wrapf(yaruserror.ErrNotSet, "RelationsValues for relation ID = %v not init", propertyID)
	}

	relationValue, ok := e.RelationsValues[propertyID]
	if !ok {
		return errors.Wrapf(yaruserror.ErrNotFound, "RelationsValues for relation ID = %v not found", propertyID)
	}

	intValues := make([]int, 0, len(values))
	for _, value := range values {
		intValue, err := property.GetValueInt(value)
		if err != nil {
			return err
		}
		intValues = append(intValues, intValue)
	}

	err := relationValue.RemoveValues(values)
	if err != nil {
		if !errors.As(err, yaruserror.ErrNotFoundList{}) {
			return err
		}
		errList, ok := err.(yaruserror.ErrNotFoundList)
		if ok {
			notFound = errList.List()
		}
	}

	if len(relationValue.Value) == 0 {
		e.DeletePropertyValues(propertyID)
	} else {
		e.setRelationValue(&relationValue)

		foundedIntValues := make([]int, 0, len(intValues)-len(notFound))
		for i, intValue := range intValues {
			if _, ok := notFound[i]; !ok {
				foundedIntValues = append(foundedIntValues, intValue)
			}
		}
		sort.Ints(foundedIntValues)

		for i := range e.IntValues {
			if e.IntValues[i].PropertyID == propertyID {
				idx := sort.SearchInts(foundedIntValues, e.IntValues[i].Value)
				// if exists
				if idx < len(foundedIntValues) && foundedIntValues[idx] == e.IntValues[i].Value {
					e.IntValues = append(e.IntValues[:i], e.IntValues[i+1:]...)
				}
			}
		}
	}

	return err
}

// По заданному значению relation привязываем relatedEntityID
func (e *Entity) BindRelatedEntityID(relation *entity_type.Relation, relatedEntityID uint) error {
	return e.AddValueForRelation(relation, relatedEntityID)
}

// По заданному значению relation привязываем relatedEntityIDs
func (e *Entity) BindRelatedEntityIDs(relation *entity_type.Relation, relatedEntityIDs []uint) error {
	return e.AddValuesForRelation(relation, relatedEntityIDs)
}

// По заданному значению relation отвязываем relatedEntityID
func (e *Entity) UnbindRelatedEntityID(relation *entity_type.Relation, relatedEntityID uint) error {
	return e.RemoveValueForRelation(relation, relatedEntityID)
}

// По заданному значению relation отвязываем relatedEntityIDs
func (e *Entity) UnbindRelatedEntityIDs(relation *entity_type.Relation, relatedEntityIDs []uint) error {
	return e.RemoveValuesForRelation(relation, relatedEntityIDs)
}
