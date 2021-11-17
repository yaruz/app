package entity

import (
	"context"

	"github.com/minipkg/selection_condition"
)

type SearchService interface {
	Get(ctx context.Context, id uint, langID uint) (*Entity, error)
	First(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) (*Entity, error)
	Query(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]Entity, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error)
}
