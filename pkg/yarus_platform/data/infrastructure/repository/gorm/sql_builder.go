package gorm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/minipkg/selection_condition"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/bool_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/date_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/float_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/int_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/time_value"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
)

type SearchResult struct {
	ID            uint
	EntityTypeID  uint
	BID           uint
	BPropertyID   uint
	BValue        bool
	IID           uint
	IPropertyID   uint
	IValue        int
	FID           uint
	FPropertyID   uint
	FValue        float64
	DID           uint
	DPropertyID   uint
	DValue        time.Time
	TID           uint
	TPropertyID   uint
	TValue        time.Time
	TxtID         uint
	TxtPropertyID uint
	TxtValue      string
}

type sqlBuilder struct {
	PropertyFinder entity.PropertyFinder
	Condition      *SearchCondition
	LangID         uint
	From           []string
	Where          *whereBuilder
	SortOrder      []string
	Limit          uint
	Offset         uint
}

type whereBuilder struct {
	Str    []string
	Params []interface{}
}

func (b *sqlBuilder) Process() {
	b.ProcessEntityConditionWhere()
	b.ProcessEntityConditionSortOrder()
	b.ProcessEntityConditionLimitOffset()
	b.ProcessPropertiesConditions()
}

func (b *sqlBuilder) ProcessEntityConditionWhere() {
	b.From = append(b.From, " "+entity.TableName+" ")

	if b.Condition.EntityCondition.Where == nil || len(b.Condition.EntityCondition.Where) == 0 {
		return
	}

	b.initWhere(len(b.Condition.EntityCondition.Where))

	for _, wc := range b.Condition.EntityCondition.Where {
		var fieldName string

		switch wc.Field {
		case fieldName_ID:
			fieldName = "id"
		case fieldName_EntityTypeID:
			fieldName = "entity_type_id"
		default:
			continue
		}

		switch wc.Condition {
		case selection_condition.ConditionEq:
			b.Where.Str = append(b.Where.Str, entity.TableName+"."+fieldName+" = ?")
			b.Where.Params = append(b.Where.Params, wc.Value)
		case selection_condition.ConditionIn:
			b.Where.Str = append(b.Where.Str, entity.TableName+"."+fieldName+" in (?)")
			b.Where.Params = append(b.Where.Params, wc.Value)
		}
	}
}

func (b *sqlBuilder) ProcessEntityConditionSortOrder() {
	if b.Condition.EntityCondition.SortOrder == nil || len(b.Condition.EntityCondition.SortOrder) == 0 {
		return
	}

	b.initSortOrder(len(b.Condition.EntityCondition.SortOrder))

	for _, sortOrderMap := range b.Condition.EntityCondition.SortOrder {
		for field, sortOrder := range sortOrderMap {
			var fieldName string

			switch field {
			case fieldName_ID:
				fieldName = "id"
			case fieldName_EntityTypeID:
				fieldName = "entity_type_id"
			default:
				continue
			}

			b.SortOrder = append(b.SortOrder, fieldName+" "+sortOrder)
		}
	}
}

func (b *sqlBuilder) ProcessEntityConditionLimitOffset() {
	b.Limit = b.Condition.EntityCondition.Limit
	b.Offset = b.Condition.EntityCondition.Offset
}

func (b *sqlBuilder) initWhere(cap int) {
	if b.Where == nil {
		b.Where = &whereBuilder{
			Str:    make([]string, 0, cap),
			Params: make([]interface{}, 0, cap),
		}
	}
}

func (b *sqlBuilder) initSortOrder(cap int) {
	if b.SortOrder == nil {
		b.SortOrder = make([]string, 0, cap)
	}
}

func (b *sqlBuilder) ProcessPropertiesConditions() error {
	ctx := context.Background()
	for propertyID, propertyCondition := range b.Condition.PropertiesConditions {
		typeID, err := b.PropertyFinder.GetTypeIDByID(ctx, propertyID)
		if err != nil {
			return err
		}

		tableAlias := fmt.Sprintf("t%vi%v", typeID, propertyID)
		b.From = append(b.From, b.JoinPropertyValue(b.GetPropertyValueTable(typeID), tableAlias, propertyID))

		b.ProcessPropertyConditionsWhere(tableAlias, propertyCondition.Where)
		b.ProcessPropertyConditionsSortOrder(tableAlias, propertyCondition.SortOrder)
	}
	return nil
}

func (b *sqlBuilder) ProcessPropertyConditionsWhere(tableAlias string, wcs selection_condition.WhereConditions) error {
	if wcs == nil || len(wcs) == 0 {
		return nil
	}

	b.initWhere(len(wcs))

	for _, wc := range wcs {
		switch wc.Condition {
		case selection_condition.ConditionEq:
			b.Where.Str = append(b.Where.Str, tableAlias+".value = ?")
			b.Where.Params = append(b.Where.Params, wc.Value)
		case selection_condition.ConditionGt:
			b.Where.Str = append(b.Where.Str, tableAlias+".value > ?")
			b.Where.Params = append(b.Where.Params, wc.Value)
		case selection_condition.ConditionGte:
			b.Where.Str = append(b.Where.Str, tableAlias+".value >= ?")
			b.Where.Params = append(b.Where.Params, wc.Value)
		case selection_condition.ConditionLt:
			b.Where.Str = append(b.Where.Str, tableAlias+".value < ?")
			b.Where.Params = append(b.Where.Params, wc.Value)
		case selection_condition.ConditionLte:
			b.Where.Str = append(b.Where.Str, tableAlias+".value <= ?")
			b.Where.Params = append(b.Where.Params, wc.Value)
		case selection_condition.ConditionIn:
			value, ok := wc.Value.([]interface{})
			if !ok {
				return errors.Errorf("Can not cast into a slice value = %v", wc.Value)
			}
			b.Where.Str = append(b.Where.Str, tableAlias+".value IN (?)")
			b.Where.Params = append(b.Where.Params, value)
		case selection_condition.ConditionBt:
			value, ok := wc.Value.([]interface{})
			if !ok {
				return errors.Errorf("Can not cast into a slice value = %v", wc.Value)
			}
			if len(value) != 2 {
				return errors.Errorf("Length of a slice must be = 2, %v given.", wc.Value)
			}
			b.Where.Str = append(b.Where.Str, tableAlias+".value BEETWIN ? AND ?")
			b.Where.Params = append(b.Where.Params, value[0], value[1])
		}

	}
	return nil
}

func (b *sqlBuilder) ProcessPropertyConditionsSortOrder(tableAlias string, sortOrderMaps []map[string]string) {
	if sortOrderMaps == nil || len(sortOrderMaps) == 0 {
		return
	}

	b.initSortOrder(len(sortOrderMaps))

	for _, sortOrderMap := range sortOrderMaps {
		for _, sortOrder := range sortOrderMap {
			b.SortOrder = append(b.SortOrder, tableAlias+".value "+sortOrder)
		}
	}
}

func (b *sqlBuilder) GetPropertyValueTable(propertyTypeID uint) string {
	var tableName string
	switch propertyTypeID {
	case property_type.IDRelation, property_type.IDInt:
		tableName = int_value.TableName
	case property_type.IDBoolean:
		tableName = bool_value.TableName
	case property_type.IDFloat:
		tableName = float_value.TableName
	case property_type.IDDate:
		tableName = date_value.TableName
	case property_type.IDTime:
		tableName = time_value.TableName
	case property_type.IDText:
		tableName = text_value.TableName
	}
	return tableName
}

func (b *sqlBuilder) JoinPropertyValue(tableName string, tableAlias string, propertyID uint) string {
	return fmt.Sprintf("INNER JOIN %v AS %v ON %v.id = %v.entity_id AND %v.property_id = %v", tableName, tableAlias, entity.TableName, tableAlias, tableAlias, propertyID)
}

func (b *sqlBuilder) subquery4Select(limit uint) (string, []interface{}) {
	strBuilder := strings.Builder{}
	strBuilder.WriteString("SELECT entity.id, row_number() OVER(ORDER BY " + strings.Join(b.SortOrder, ", ") + ")")
	strBuilder.WriteString(" FROM " + strings.Join(b.From, " "))
	strBuilder.WriteString(" WHERE " + strings.Join(b.Where.Str, " AND "))

	if limit > 0 {
		strBuilder.WriteString(fmt.Sprintf(" Limit %v", limit))
	}

	if b.Offset > 0 {
		strBuilder.WriteString(fmt.Sprintf(" OFFSET %v", b.Offset))
	}

	return strBuilder.String(), b.Where.Params
}

func (b *sqlBuilder) CountQuery() (string, []interface{}) {
	strBuilder := strings.Builder{}
	strBuilder.WriteString("SELECT COUNT(entity.id)")
	strBuilder.WriteString(" FROM " + strings.Join(b.From, " "))
	strBuilder.WriteString(" WHERE " + strings.Join(b.Where.Str, " AND "))
	strBuilder.WriteString(" ORDER BY " + strings.Join(b.SortOrder, ", "))

	return strBuilder.String(), b.Where.Params
}

func (b *sqlBuilder) mainPartOfQuery(subQuery string, subQueryParams []interface{}) (string, []interface{}) {
	tpl := `select e.*, b.id as b_id, b.property_id as b_property_id, b.value as b_value,
i.id as i_id, i.property_id as i_property_id, i.value as i_value,
f.id as f_id, f.property_id as f_property_id, f.value as f_value,
d.id as d_id, d.property_id as d_property_id, d.value as d_value,
t.id as t_id, t.property_id as t_property_id, t.value as t_value,
txt.id as txt_id, txt.property_id as txt_property_id, txt.value as txt_value
from entity e 
inner join (%s) as x(id, sort_order) on e.id = x.id
left join bool_value b on e.id = b.entity_id 
left join int_value i on e.id = i.entity_id 
left join float_value f on e.id = f.entity_id 
left join date_value d on e.id = d.entity_id 
left join time_value t on e.id = t.entity_id 
left join text_value txt on e.id = txt.entity_id and txt.lang_id = ?
order by x.sort_order`
	return fmt.Sprintf(tpl, subQuery), append(subQueryParams, b.LangID)
}

func (b *sqlBuilder) SelectQuery() (string, []interface{}) {
	return b.mainPartOfQuery(b.subquery4Select(b.Limit))
}

func (b *sqlBuilder) FirstQuery() (string, []interface{}) {
	return b.mainPartOfQuery(b.subquery4Select(1))
}
