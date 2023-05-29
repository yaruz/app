package gorm

import (
	"context"

	minipkg_gorm "github.com/minipkg/db/gorm"

	"github.com/minipkg/log"
	"github.com/minipkg/selection_condition"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/bool_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/date_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/float_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/int_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/time_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/utext_value"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
)

type SearchService struct {
	mapReducer       IMapReducer
	logger           log.Logger
	propertyFinder   entity.PropertyFinder
	entityTypeFinder entity.EntityTypeFinder
	langFinder       entity.LangFinder
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

var _ entity.SearchService = (*SearchService)(nil)

var IDConditionVariants = []interface{}{
	selection_condition.ConditionEq,
	selection_condition.ConditionIn,
}

func NewSearchService(logger log.Logger, mapReducer IMapReducer, entityTypeFinder entity.EntityTypeFinder, propertyFinder entity.PropertyFinder, langFinder entity.LangFinder) *SearchService {
	return &SearchService{
		logger:           logger,
		mapReducer:       mapReducer,
		propertyFinder:   propertyFinder,
		entityTypeFinder: entityTypeFinder,
		langFinder:       langFinder,
		model:            entity.New(),
	}
}

func (s *SearchService) Get(ctx context.Context, ID uint, typeID uint, langID uint) (*entity.Entity, error) {
	sqlBuilder := &sqlBuilder{
		Ctx:    context.Background(),
		LangID: langID,
	}
	sql, params := sqlBuilder.GetQuery(ID)

	db, err := s.mapReducer.GetDB(ctx, typeID, ID)
	if err != nil {
		return nil, err
	}

	searchResult := make([]SearchResult, 0)
	if err := db.DB().Raw(sql, params...).Scan(&searchResult).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, yaruserror.ErrNotFound
		}
		return nil, err
	}

	return s.instantiateItem(searchResult), nil
}

func (s *SearchService) First(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (*entity.Entity, error) {
	parser, searchCondition, err := s.ParseSelectionCondition(condition)
	if err != nil {
		return nil, err
	}

	sqlBuilder := s.newSqlBuilder(searchCondition, langID)
	sql, params := sqlBuilder.FirstQuery()

	searchResults, err := s.mapReducer.Query(ctx, parser, searchCondition.EntityCondition.Where, func(db minipkg_gorm.IDB) ([]SearchResult, error) {
		searchResult := make([]SearchResult, 0)
		if err = db.DB().Raw(sql, params...).Scan(&searchResult).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, yaruserror.ErrNotFound
			}
			return nil, err
		}
		return searchResult, nil
	})
	if err != nil {
		return nil, err
	}

	return s.instantiateItem(searchResults), nil
}

func (s *SearchService) Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]entity.Entity, error) {
	parser, searchCondition, err := s.ParseSelectionCondition(condition)
	if err != nil {
		return nil, err
	}

	sqlBuilder := s.newSqlBuilder(searchCondition, langID)
	sql, params := sqlBuilder.SelectQuery()

	searchResults, err := s.mapReducer.Query(ctx, parser, searchCondition.EntityCondition.Where, func(db minipkg_gorm.IDB) ([]SearchResult, error) {
		searchResult := make([]SearchResult, 0)
		if err = db.DB().Raw(sql, params...).Scan(&searchResult).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, yaruserror.ErrNotFound
			}
			return nil, err
		}
		return searchResult, nil
	})

	return s.instantiateItems(searchResults), nil
}

func (s *SearchService) Count(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (uint, error) {
	parser, searchCondition, err := s.ParseSelectionCondition(condition)
	if err != nil {
		return 0, err
	}

	sqlBuilder := s.newSqlBuilder(searchCondition, langID)
	sql, params := sqlBuilder.CountQuery()

	searchResults, err := s.mapReducer.Count(ctx, parser, searchCondition.EntityCondition.Where, func(db minipkg_gorm.IDB) (uint, error) {
		var searchResult uint
		if err = db.DB().Raw(sql, params...).Scan(&searchResult).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, yaruserror.ErrNotFound
			}
			return 0, err
		}
		return searchResult, nil
	})

	return searchResults, nil
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
			e = entity.New()
			e.ID = searchResult.ID
			e.EntityTypeID = searchResult.EntityTypeID
			IDs = append(IDs, e.ID)
			propValsMap[e.ID] = s.initPropValsMapItem(cap)
		}

		s.processBoolValue(e, &propValsMap, &searchResult)
		s.processIntValue(e, &propValsMap, &searchResult)
		s.processFloatValue(e, &propValsMap, &searchResult)
		s.processDateValue(e, &propValsMap, &searchResult)
		s.processTimeValue(e, &propValsMap, &searchResult)
		s.processTextValue(e, &propValsMap, &searchResult)
		s.processUTextValue(e, &propValsMap, &searchResult)

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
				e.BoolValues = make([]bool_value.BoolValue, 0, 1)
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
				e.IntValues = make([]int_value.IntValue, 0, 1)
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
				e.FloatValues = make([]float_value.FloatValue, 0, 1)
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
				e.DateValues = make([]date_value.DateValue, 0, 1)
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
				e.TimeValues = make([]time_value.TimeValue, 0, 1)
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
				e.TextValues = make([]text_value.TextValue, 0, 1)
			}

			e.TextValues = append(e.TextValues, text_value.TextValue{
				ID:         searchResult.TxtID,
				EntityID:   e.ID,
				LangID:     searchResult.TxtLangID,
				PropertyID: searchResult.TxtPropertyID,
				Value:      searchResult.TxtValue,
			})
		}
	}
}

func (s *SearchService) processUTextValue(e *entity.Entity, propValsMap *map[uint]map[uint]map[uint]struct{}, searchResult *SearchResult) {
	if searchResult.UTxtID > 0 {
		if _, ok := (*propValsMap)[e.ID][property_type.IDUText][searchResult.UTxtID]; !ok {
			(*propValsMap)[e.ID][property_type.IDUText][searchResult.UTxtID] = struct{}{}

			if e.UTextValues == nil {
				e.UTextValues = make([]utext_value.UTextValue, 0, 1)
			}

			e.UTextValues = append(e.UTextValues, utext_value.UTextValue{
				ID:         searchResult.UTxtID,
				EntityID:   e.ID,
				PropertyID: searchResult.UTxtPropertyID,
				Value:      searchResult.UTxtValue,
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

func (s *SearchService) ParseSelectionCondition(OriginalCondition *selection_condition.SelectionCondition) (*SelectionConditionParser, *SearchCondition, error) {
	parser := &SelectionConditionParser{
		propertyFinder:    s.propertyFinder,
		entityTypeFinder:  s.entityTypeFinder,
		OriginalCondition: OriginalCondition,
		TargetCondition:   &SearchCondition{},
	}

	if err := parser.Run(); err != nil {
		return nil, nil, err
	}

	return parser, parser.TargetCondition, nil
}

func (s *SearchService) newSqlBuilder(condition *SearchCondition, langID uint) *sqlBuilder {
	builder := &sqlBuilder{
		Ctx:            context.Background(),
		PropertyFinder: s.propertyFinder,
		LangFinder:     s.langFinder,
		Condition:      condition,
		LangID:         langID,
		From:           make([]string, 0, 1),
	}
	builder.Process()
	return builder
}
