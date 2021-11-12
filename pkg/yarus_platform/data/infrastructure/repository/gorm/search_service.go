package gorm

import (
	"context"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/log"
	"github.com/minipkg/selection_condition"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
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
	var searchResult SearchResult

	searchCondition, err := s.ParseSelectionCondition(condition)
	if err != nil {
		return nil, err
	}

	sqlBuilder := s.newSqlBuilder(searchCondition, langID)

	sql, params := sqlBuilder.FirstQuery()
	s.db.DB().Raw(sql, params...).Scan(&searchResult)

	return s.instantiateItem(&searchResult), nil
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

func (s *SearchService) instantiateItem(searchResult *SearchResult) *entity.Entity {
	e := entity.New()
	return e
}

func (s *SearchService) instantiateItems(searchResults []SearchResult) []entity.Entity {
	entities := make([]entity.Entity, 0, len(searchResults))
	return entities
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
