package gorm

import (
	"context"

	minipkg_gorm "github.com/minipkg/db/gorm"

	"github.com/minipkg/log"
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

type SearchService struct {
	db               minipkg_gorm.IDB
	logger           log.ILogger
	propertyFinder   entity.PropertyFinder
	entityTypeFinder entity.EntityTypeFinder
	model            *entity.Entity
}

type SearchCondition struct {
	EntityCondition      EntityCondition
	PropertiesConditions map[uint]PropertyCondition
}

type EntityCondition struct {
	Where     selection_condition.WhereConditions
	SortOrder []map[string]string
	Limit     uint
	Offset    uint
}

type PropertyCondition struct {
	Where     selection_condition.WhereConditions
	SortOrder []map[string]string
}

const (
	fieldName_ID           = "ID"
	fieldName_EntityTypeID = "EntityTypeID"
	fieldName_EntityType   = "EntityType"
)

var _ entity.SearchService = (*SearchService)(nil)

var IDConditionVariants = []interface{}{
	selection_condition.ConditionEq,
	selection_condition.ConditionIn,
}

func NewSearchService(logger log.ILogger, dbase minipkg_gorm.IDB, entityTypeFinder entity.EntityTypeFinder, propertyFinder entity.PropertyFinder) (*SearchService, error) {
	var err error
	ctx := context.Background()
	service := &SearchService{
		logger:           logger,
		propertyFinder:   propertyFinder,
		entityTypeFinder: entityTypeFinder,
		model:            entity.New(),
	}
	if service.db, err = dbase.SchemeInitWithContext(ctx, service.model); err != nil {
		return nil, err
	}
	return service, nil
}

func (s *SearchService) First(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (*entity.Entity, error) {
	searchCondition, err := s.ParseSelectionCondition(condition)
	if err != nil {
		return nil, err
	}

	sqlBuilder := s.newSqlBuilder(searchCondition, langID)

	searchResults := make([]SearchResult, 0)
	sql, params := sqlBuilder.FirstQuery()
	s.db.DB().Raw(sql, params...).Scan(&searchResults)

	return s.instantiateItem(searchResults), nil
}

func (s *SearchService) Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]entity.Entity, error) {
	searchCondition, err := s.ParseSelectionCondition(condition)
	if err != nil {
		return nil, err
	}

	sqlBuilder := s.newSqlBuilder(searchCondition, langID)

	searchResults := make([]SearchResult, 0)
	sql, params := sqlBuilder.SelectQuery()

	err = s.db.DB().
		Raw(sql, params...).
		Scan(&searchResults).
		Error

	if err != nil {
		return nil, err
	}

	return s.instantiateItems(searchResults), nil
}

func (s *SearchService) Count(ctx context.Context, condition *selection_condition.SelectionCondition) (uint, error) {
	searchCondition, err := s.ParseSelectionCondition(condition)
	if err != nil {
		return 0, err
	}

	sqlBuilder := s.newSqlBuilder(searchCondition, 0)

	var res uint
	sql, params := sqlBuilder.CountQuery()
	s.db.DB().Raw(sql, params...).Scan(&res)
	return res, nil
}

func (s *SearchService) instantiateItem(searchResults []SearchResult) *entity.Entity {
	if searchResults == nil || len(searchResults) == 0 {
		return nil
	}
	results := s.instantiateItems(searchResults)

	if results == nil || len(results) == 0 {
		return nil
	}
	return &results[0]
}

func (s *SearchService) instantiateItems(searchResults []SearchResult) []entity.Entity {
	cap := len(searchResults) // небольшой оверхед, чтобы не было аллокаций
	IDs := make([]uint, 0, cap)
	entitiesMap := make(map[uint]*entity.Entity, cap)
	propValsMap := make(map[uint]map[uint]map[uint]struct{}, cap) // entity_id -> propperty_type_id -> property_value_id

	for _, searchResult := range searchResults {
		var e *entity.Entity
		var ok bool

		if e, ok = entitiesMap[searchResult.ID]; !ok {
			e := entity.New()
			e.ID = searchResult.ID
			e.EntityTypeID = searchResult.EntityTypeID
			entitiesMap[searchResult.ID] = e
			IDs = append(IDs, e.ID)
			propValsMap[e.ID] = s.initPropValsMapItem(cap)
		}

		s.processBoolValue(e, &propValsMap, &searchResult)
		s.processIntValue(e, &propValsMap, &searchResult)
		s.processFloatValue(e, &propValsMap, &searchResult)
		s.processDateValue(e, &propValsMap, &searchResult)
		s.processTimeValue(e, &propValsMap, &searchResult)
		s.processTextValue(e, &propValsMap, &searchResult)

		entitiesMap[e.ID] = e
	}

	entities := make([]entity.Entity, 0, len(IDs))
	for _, id := range IDs {
		entities = append(entities, *(entitiesMap[id]))
	}

	return entities
}

func (s *SearchService) processBoolValue(e *entity.Entity, propValsMap *map[uint]map[uint]map[uint]struct{}, searchResult *SearchResult) {
	if searchResult.BID > 0 {
		if _, ok := (*propValsMap)[e.ID][property_type.IDBoolean][searchResult.BID]; !ok {
			(*propValsMap)[e.ID][property_type.IDBoolean][searchResult.BID] = struct{}{}

			if e.BoolValues == nil {
				e.BoolValues = make([]bool_value.BoolValue, 1)
			}

			e.BoolValues = append(e.BoolValues, bool_value.BoolValue{
				ID:         searchResult.BID,
				EntityID:   e.ID,
				PropertyID: searchResult.BPropertyID,
				Value:      searchResult.BValue,
			})
		}
	}
}

func (s *SearchService) processIntValue(e *entity.Entity, propValsMap *map[uint]map[uint]map[uint]struct{}, searchResult *SearchResult) {
	if searchResult.IID > 0 {
		if _, ok := (*propValsMap)[e.ID][property_type.IDInt][searchResult.IID]; !ok {
			(*propValsMap)[e.ID][property_type.IDInt][searchResult.IID] = struct{}{}

			if e.IntValues == nil {
				e.IntValues = make([]int_value.IntValue, 1)
			}

			e.IntValues = append(e.IntValues, int_value.IntValue{
				ID:         searchResult.IID,
				EntityID:   e.ID,
				PropertyID: searchResult.IPropertyID,
				Value:      searchResult.IValue,
			})
		}
	}
}

func (s *SearchService) processFloatValue(e *entity.Entity, propValsMap *map[uint]map[uint]map[uint]struct{}, searchResult *SearchResult) {
	if searchResult.FID > 0 {
		if _, ok := (*propValsMap)[e.ID][property_type.IDFloat][searchResult.FID]; !ok {
			(*propValsMap)[e.ID][property_type.IDFloat][searchResult.FID] = struct{}{}

			if e.FloatValues == nil {
				e.FloatValues = make([]float_value.FloatValue, 1)
			}

			e.FloatValues = append(e.FloatValues, float_value.FloatValue{
				ID:         searchResult.FID,
				EntityID:   e.ID,
				PropertyID: searchResult.FPropertyID,
				Value:      searchResult.FValue,
			})
		}
	}
}

func (s *SearchService) processDateValue(e *entity.Entity, propValsMap *map[uint]map[uint]map[uint]struct{}, searchResult *SearchResult) {
	if searchResult.DID > 0 {
		if _, ok := (*propValsMap)[e.ID][property_type.IDDate][searchResult.DID]; !ok {
			(*propValsMap)[e.ID][property_type.IDDate][searchResult.DID] = struct{}{}

			if e.DateValues == nil {
				e.DateValues = make([]date_value.DateValue, 1)
			}

			e.DateValues = append(e.DateValues, date_value.DateValue{
				ID:         searchResult.DID,
				EntityID:   e.ID,
				PropertyID: searchResult.DPropertyID,
				Value:      searchResult.DValue,
			})
		}
	}
}

func (s *SearchService) processTimeValue(e *entity.Entity, propValsMap *map[uint]map[uint]map[uint]struct{}, searchResult *SearchResult) {
	if searchResult.TID > 0 {
		if _, ok := (*propValsMap)[e.ID][property_type.IDTime][searchResult.TID]; !ok {
			(*propValsMap)[e.ID][property_type.IDTime][searchResult.TID] = struct{}{}

			if e.TimeValues == nil {
				e.TimeValues = make([]time_value.TimeValue, 1)
			}

			e.TimeValues = append(e.TimeValues, time_value.TimeValue{
				ID:         searchResult.TID,
				EntityID:   e.ID,
				PropertyID: searchResult.TPropertyID,
				Value:      searchResult.TValue,
			})
		}
	}
}

func (s *SearchService) processTextValue(e *entity.Entity, propValsMap *map[uint]map[uint]map[uint]struct{}, searchResult *SearchResult) {
	if searchResult.TxtID > 0 {
		if _, ok := (*propValsMap)[e.ID][property_type.IDText][searchResult.TxtID]; !ok {
			(*propValsMap)[e.ID][property_type.IDText][searchResult.TxtID] = struct{}{}

			if e.TextValues == nil {
				e.TextValues = make([]text_value.TextValue, 1)
			}

			e.TextValues = append(e.TextValues, text_value.TextValue{
				ID:         searchResult.TxtID,
				EntityID:   e.ID,
				PropertyID: searchResult.TxtPropertyID,
				Value:      searchResult.TxtValue,
			})
		}
	}
}

func (s *SearchService) initPropValsMapItem(cap int) map[uint]map[uint]struct{} {
	res := make(map[uint]map[uint]struct{}, property_type.Count)
	for i := uint(property_type.IDBoolean); i < property_type.Count+property_type.IDBoolean; i++ {
		res[i] = make(map[uint]struct{}, cap)
	}
	return res
}

func (s *SearchService) ParseSelectionCondition(OriginalCondition *selection_condition.SelectionCondition) (*SearchCondition, error) {
	parser := SelectionConditionParser{
		propertyFinder:    s.propertyFinder,
		entityTypeFinder:  s.entityTypeFinder,
		OriginalCondition: OriginalCondition,
		TargetCondition:   &SearchCondition{},
	}

	if err := parser.Run(); err != nil {
		return nil, err
	}

	return parser.TargetCondition, nil
}

func (s *SearchService) newSqlBuilder(condition *SearchCondition, langID uint) *sqlBuilder {
	builder := &sqlBuilder{
		PropertyFinder: s.propertyFinder,
		Condition:      condition,
		LangID:         langID,
		From:           make([]string, 0, 1),
	}
	builder.Process()
	return builder
}
