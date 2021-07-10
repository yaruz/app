package entity_type

import (
	"context"

	"github.com/minipkg/selection_condition"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
)

// RelationService encapsulates usecase logic.
type RelationService interface {
	NewEntity() *Relation
	Get(ctx context.Context, id uint) (*Relation, error)
	Query(ctx context.Context, query *selection_condition.SelectionCondition) ([]Relation, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *Relation) error
	Update(ctx context.Context, entity *Relation) error
	Save(ctx context.Context, entity *Relation) error
	Delete(ctx context.Context, id uint) error
	PropertyAndRelationQuery(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property.Property, []Relation, error)
	GetPropertiesAndRelationsByEntityTypeID(ctx context.Context, entityTypeID uint) ([]property.Property, []Relation, error)
}

type relationService struct {
	logger     log.ILogger
	repository RelationRepository
}

var _ RelationService = (*relationService)(nil)

// NewRelationService creates a new relationService.
func NewRelationService(logger log.ILogger, repo RelationRepository) RelationService {
	s := &relationService{
		logger:     logger,
		repository: repo,
	}
	repo.SetDefaultConditions(s.defaultConditions())
	return s
}

// Defaults returns defaults params
func (s *relationService) defaultConditions() *selection_condition.SelectionCondition {
	return &selection_condition.SelectionCondition{}
}

func (s *relationService) NewEntity() *Relation {
	return NewRelation()
}

// Get returns the entity with the specified ID.
func (s *relationService) Get(ctx context.Context, id uint) (*Relation, error) {
	entity, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Query returns the items with the specified selection condition.
func (s *relationService) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]Relation, error) {
	items, err := s.repository.Query(ctx, cond)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of items by query: %v", cond)
	}
	return items, nil
}

func (s *relationService) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	count, err := s.repository.Count(ctx, cond)
	if err != nil {
		return 0, errors.Wrapf(err, "Can not count a list of items by query: %v", cond)
	}
	return count, nil
}

func (s *relationService) Create(ctx context.Context, entity *Relation) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.Create(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not create an entity: %v", entity)
	}
	return nil
}

func (s *relationService) Update(ctx context.Context, entity *Relation) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.Update(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not update an entity: %v", entity)
	}
	return nil
}

func (s *relationService) Save(ctx context.Context, entity *Relation) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.Save(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not save an entity: %v", entity)
	}
	return nil
}

func (s *relationService) Delete(ctx context.Context, id uint) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "Can not delete an entity by ID: %v", id)
	}
	return nil
}

func (s *relationService) PropertyAndRelationQuery(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property.Property, []Relation, error) {
	props, rels, err := s.repository.PropertyAndRelationQuery(ctx, cond)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Can not get properties and relations by query: %v", cond)
	}
	return props, rels, nil
}

func (s *relationService) GetPropertiesAndRelationsByEntityTypeID(ctx context.Context, entityTypeID uint) ([]property.Property, []Relation, error) {
	props, rels, err := s.repository.GetPropertiesAndRelationsByEntityTypeID(ctx, entityTypeID)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Can not get properties and relations by entityTypeID: %v", entityTypeID)
	}
	return props, rels, nil
}
