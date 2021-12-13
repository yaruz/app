package gorm

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/minipkg/selection_condition"
	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/sysname"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
)

type SelectionConditionParser struct {
	propertyFinder    entity.PropertyFinder
	entityTypeFinder  entity.EntityTypeFinder
	OriginalCondition *selection_condition.SelectionCondition
	TargetCondition   *SearchCondition
}

func (p *SelectionConditionParser) Run() error {
	if err := p.ParseWhere(); err != nil {
		return err
	}

	if err := p.ParseSortOrder(); err != nil {
		return err
	}

	if err := p.ParseLimitOffset(); err != nil {
		return err
	}
	return nil
}

func (p *SelectionConditionParser) ParseWhere() error {
	ctx := context.Background()
	sourceWhereConditions := make(selection_condition.WhereConditions, 0)

	if p.OriginalCondition.Where != nil {
		switch wc := p.OriginalCondition.Where.(type) {
		case selection_condition.WhereConditions:
			sourceWhereConditions = append(sourceWhereConditions, wc...)
		case selection_condition.WhereCondition:
			sourceWhereConditions = append(sourceWhereConditions, wc)
		default:
			return errors.Errorf("Conditions must be only a selection_condition.WhereConditions or a selection_condition.WhereCondition. ")
		}
	}

	for _, wcondition := range sourceWhereConditions {
		if wcondition.Field == entity.FieldName_ID || wcondition.Field == entity.FieldName_EntityTypeID {
			wc, err := p.CheckIDCondition(&wcondition)
			if err != nil {
				return err
			}
			p.initEntityConditionWhere()
			p.TargetCondition.EntityCondition.Where = append(p.TargetCondition.EntityCondition.Where, *wc)
			continue
		}

		if wcondition.Field == entity.FieldName_EntityType {
			wc, err := p.CheckEntityTypeCondition(&wcondition)
			if err != nil {
				return err
			}
			p.initEntityConditionWhere()
			p.TargetCondition.EntityCondition.Where = append(p.TargetCondition.EntityCondition.Where, *wc)
			continue
		}
		// Нужен-ли?
		propertySysname, err := p.CheckPropertySysname(wcondition.Field)
		if err != nil {
			return err
		}

		propertyID, err := p.propertyFinder.GetIDBySysname(ctx, propertySysname)
		if err != nil {
			return err
		}

		p.initPropertiesConditions()
		propertiesConditions, ok := p.TargetCondition.PropertiesConditions[propertyID]

		if !ok {
			propertiesConditions = PropertyCondition{}
		}

		if propertiesConditions.Where == nil {
			propertiesConditions.Where = make(selection_condition.WhereConditions, 0, 1)
		}
		propertiesConditions.Where = append(propertiesConditions.Where, wcondition)
		p.TargetCondition.PropertiesConditions[propertyID] = propertiesConditions
	}

	return nil

}

func (p *SelectionConditionParser) initEntityConditionWhere() {
	if p.TargetCondition.EntityCondition.Where == nil {
		p.TargetCondition.EntityCondition.Where = make(selection_condition.WhereConditions, 0, 1)
	}
}

func (p *SelectionConditionParser) initEntityConditionSortOrder() {
	if p.TargetCondition.EntityCondition.SortOrder == nil {
		p.TargetCondition.EntityCondition.SortOrder = make([]map[string]string, 0, 1)
	}
}

func (p *SelectionConditionParser) initPropertiesConditions() {
	if p.TargetCondition.PropertiesConditions == nil {
		p.TargetCondition.PropertiesConditions = make(map[uint]PropertyCondition, 0)
	}
}

func (p *SelectionConditionParser) ParseSortOrder() error {
	ctx := context.Background()
	if err := p.CheckSortOrder(p.OriginalCondition.SortOrder); err != nil {
		return err
	}

	for _, sortOrderMap := range p.OriginalCondition.SortOrder {

		for field, sortOrder := range sortOrderMap {

			if field == entity.FieldName_ID {
				p.initEntityConditionSortOrder()
				p.TargetCondition.EntityCondition.SortOrder = append(p.TargetCondition.EntityCondition.SortOrder, map[string]string{entity.FieldName_ID: sortOrder})
				continue
			}

			if field == entity.FieldName_EntityType || field == entity.FieldName_EntityTypeID {
				p.initEntityConditionSortOrder()
				p.TargetCondition.EntityCondition.SortOrder = append(p.TargetCondition.EntityCondition.SortOrder, map[string]string{entity.FieldName_EntityTypeID: sortOrder})
				continue
			}
			propertySysname := field
			propertyID, err := p.propertyFinder.GetIDBySysname(ctx, propertySysname)
			if err != nil {
				return err
			}

			p.initPropertiesConditions()
			propertiesConditions, ok := p.TargetCondition.PropertiesConditions[propertyID]

			if !ok {
				propertiesConditions = PropertyCondition{}
			}

			if propertiesConditions.SortOrder == nil {
				propertiesConditions.SortOrder = make([]map[string]string, 0, 1)
			}

			propertiesConditions.SortOrder = append(propertiesConditions.SortOrder, map[string]string{field: sortOrder})
			p.TargetCondition.PropertiesConditions[propertyID] = propertiesConditions
		}
	}

	return nil
}

func (p *SelectionConditionParser) ParseLimitOffset() error {
	p.TargetCondition.EntityCondition.Limit = p.OriginalCondition.Limit
	p.TargetCondition.EntityCondition.Offset = p.OriginalCondition.Offset
	return nil
}

func (p *SelectionConditionParser) CheckSortOrder(value []map[string]string) error {
	for _, orderMap := range value {
		for field, sortOrder := range orderMap {
			if field == entity.FieldName_ID || field == entity.FieldName_EntityType || field == entity.FieldName_EntityTypeID {
				continue
			}

			if _, err := p.propertyFinder.GetIDBySysname(context.Background(), field); err != nil {
				return err
			}
			if err := validation.Validate(sortOrder, validation.In(selection_condition.SortOrderVariants...)); err != nil {
				return errors.Wrapf(err, "Invalid SortOrder = %v", value)
			}
		}
	}
	return nil
}

func (p *SelectionConditionParser) GetEntityTypeIDsAndIDsByEntityWhereConditions(entityWhereConditions selection_condition.WhereConditions) (entityTypes []uint, IDs []uint, err error) {
	for _, wc := range entityWhereConditions {
		tids, ids, err := p.getEntityTypeIDsAndIDsByEntityWhereCondition(&wc)
		if err != nil {
			return nil, nil, err
		}
		entityTypes = append(entityTypes, tids...)
		IDs = append(IDs, ids...)
	}
	return entityTypes, IDs, nil
}

func (p *SelectionConditionParser) getEntityTypeIDsAndIDsByEntityWhereCondition(entityWhereCondition *selection_condition.WhereCondition) (entityTypes []uint, IDs []uint, err error) {
	if entityWhereCondition.Field != entity.FieldName_ID && entityWhereCondition.Field != entity.FieldName_EntityTypeID {
		return nil, nil, errors.Errorf("Field must be %q or %q.", entity.FieldName_ID, entity.FieldName_EntityTypeID)
	}
	var values []uint

	switch entityWhereCondition.Condition {
	case selection_condition.ConditionEq:
		var value uint
		if value, err = p.CheckID(entityWhereCondition.Value); err != nil {
			return nil, nil, err
		}
		values = append(values, value)
	case selection_condition.ConditionIn:
		var value []uint
		if value, err = p.CheckIDs(entityWhereCondition.Value); err != nil {
			return nil, nil, err
		}
		values = append(values, value...)
	default:
		return nil, nil, errors.Errorf("Condition must be %q or %q", selection_condition.ConditionEq, selection_condition.ConditionIn)
	}

	if entityWhereCondition.Field == entity.FieldName_ID {
		return nil, values, nil
	}

	return values, nil, nil
}

func (p *SelectionConditionParser) CheckIDCondition(wcondition *selection_condition.WhereCondition) (*selection_condition.WhereCondition, error) {
	if wcondition.Field != entity.FieldName_ID && wcondition.Field != entity.FieldName_EntityTypeID {
		return nil, errors.Errorf("Field must be %q or %q.", entity.FieldName_ID, entity.FieldName_EntityTypeID)
	}
	wc := &selection_condition.WhereCondition{
		Field:     wcondition.Field,
		Condition: wcondition.Condition,
	}
	var err error

	switch wcondition.Condition {
	case selection_condition.ConditionEq:
		var value uint
		if value, err = p.CheckID(wcondition.Value); err != nil {
			return nil, err
		}
		wc.Value = value
	case selection_condition.ConditionIn:
		var value []uint
		if value, err = p.CheckIDs(wcondition.Value); err != nil {
			return nil, err
		}
		wc.Value = value
	default:
		return nil, errors.Errorf("Condition must be %q or %q", selection_condition.ConditionEq, selection_condition.ConditionIn)
	}

	return wc, nil
}

func (p *SelectionConditionParser) CheckID(value interface{}) (uint, error) {
	return property.GetRelationItemValue(value)
}

func (p *SelectionConditionParser) CheckIDs(value interface{}) ([]uint, error) {
	valInterface, ok := value.([]interface{})
	if !ok {
		return nil, errors.Errorf("Can not cast to the []interface{} value = %v", value)
	}

	values := make([]uint, len(valInterface))
	for i, val := range valInterface {
		valUint, err := p.CheckID(val)
		if err != nil {
			return nil, err
		}
		values[i] = valUint
	}
	return values, nil
}

func (p *SelectionConditionParser) CheckEntityTypeCondition(wcondition *selection_condition.WhereCondition) (*selection_condition.WhereCondition, error) {
	if wcondition.Field != entity.FieldName_EntityType {
		return nil, errors.Errorf("Field must be %q.", entity.FieldName_EntityType)
	}
	wc := &selection_condition.WhereCondition{
		Field:     entity.FieldName_EntityTypeID,
		Condition: wcondition.Condition,
	}
	var err error

	switch wcondition.Condition {
	case selection_condition.ConditionEq:
		var value uint
		if value, err = p.CheckEntityTypeSysnameAndGetID(wcondition.Value); err != nil {
			return nil, err
		}
		wc.Value = value
	case selection_condition.ConditionIn:
		valInterface, ok := wcondition.Value.([]interface{})
		if !ok {
			return nil, errors.Errorf("Can not cast to the []interface{} value = %v", wcondition.Value)
		}

		value := make([]uint, len(valInterface))
		for i, val := range valInterface {
			valUint, err := p.CheckEntityTypeSysnameAndGetID(val)
			if err != nil {
				return nil, err
			}
			value[i] = valUint
		}
		wc.Value = value
	default:
		return nil, errors.Errorf("Condition must be %q or %q", selection_condition.ConditionEq, selection_condition.ConditionIn)
	}

	return wc, nil
}

func (p *SelectionConditionParser) CheckEntityTypeSysnameAndGetID(value interface{}) (uint, error) {
	val, err := property.GetValueText(value)
	if err != nil {
		return 0, err
	}

	if err := validation.Validate(val, sysname.ValidationRules...); err != nil {
		return 0, err
	}

	res, err := p.entityTypeFinder.GetIDBySysname(context.Background(), val)
	if err != nil {
		if errors.Is(err, yaruserror.ErrNotFound) {
			return 0, errors.Wrapf(yaruserror.ErrNotFound, "EntityType = %q not found.", val)
		}
		return 0, err
	}
	return res, nil
}

func (p *SelectionConditionParser) CheckPropertySysname(value string) (string, error) {
	if err := validation.Validate(value, sysname.ValidationRules...); err != nil {
		return "", err
	}

	if _, err := p.propertyFinder.GetIDBySysname(context.Background(), value); err != nil {
		if errors.Is(err, yaruserror.ErrNotFound) {
			return "", errors.Wrapf(yaruserror.ErrNotFound, "Property = %q not found.", value)
		}
		return "", err
	}
	return value, nil
}
