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

func (s *SearchService) First(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) (*entity.Entity, error) {

}

func (s *SearchService) Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]entity.Entity, error) {
	searchCondition, err := s.ParseSelectionCondition(condition)
	if err != nil {
		return nil, err
	}

	sqlBuilder := s.newSqlBuilder(searchCondition, langID)
}

func (s *SearchService) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {

}

func (s *SearchService) ParseSelectionCondition(OriginalCondition *selection_condition.SelectionCondition) (*SearchCondition, error) {
	parser := SelectionConditionParser{
		propertyFinder:    s.propertyFinder,
		OriginalCondition: OriginalCondition,
		TargetCondition:   &SearchCondition{},
	}

	if err := parser.Run(); err != nil {
		return nil, err
	}

	return parser.TargetCondition, nil
}

func (s *SearchService) newSqlBuilder(condition *SearchCondition, langID uint) *sqlBuilder {
	return &sqlBuilder{
		PropertyFinder: s.propertyFinder,
		Condition:      condition,
		LangID:         langID,
		From:           make([]string, 0, 1),
	}
}
