package entity_type

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"

	"github.com/minipkg/selection_condition"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
)

// RelationService encapsulates usecase logic.
type RelationService interface {
	NewEntity() *Relation
	Get(ctx context.Context, id uint) (*Relation, error)
	First(ctx context.Context, entity *Relation) (*Relation, error)
	Query(ctx context.Context, query *selection_condition.SelectionCondition) ([]Relation, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	UpsertBySysname(ctx context.Context, entity *Relation, langID uint) (err error)
	Create(ctx context.Context, entity *Relation) error
	Update(ctx context.Context, entity *Relation) error
	Delete(ctx context.Context, entity *Relation) error
	TGet(ctx context.Context, id uint, langID uint) (*Relation, error)
	TFirst(ctx context.Context, entity *Relation, langID uint) (*Relation, error)
	TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]Relation, error)
	TCreate(ctx context.Context, entity *Relation, langID uint) (err error)
	TUpdate(ctx context.Context, entity *Relation, langID uint) (err error)
	PropertyAndRelationQuery(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property.Property, []Relation, error)
	TPropertyAndRelationQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]property.Property, []Relation, error)
	GetPropertiesAndRelationsByEntityTypeID(ctx context.Context, entityTypeID uint) ([]property.Property, []Relation, error)
	TGetPropertiesAndRelationsByEntityTypeID(ctx context.Context, entityTypeID uint, langID uint) ([]property.Property, []Relation, error)
}

type relationService struct {
	logger          log.ILogger
	repository      RelationRepository
	propertyService property.IService
}

var _ RelationService = (*relationService)(nil)

// NewRelationService creates a new relationService.
func NewRelationService(logger log.ILogger, repo RelationRepository, propertyService property.IService) RelationService {
	s := &relationService{
		logger:          logger,
		repository:      repo,
		propertyService: propertyService,
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

func (s *relationService) First(ctx context.Context, entity *Relation) (*Relation, error) {
	entity, err := s.repository.First(ctx, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *relationService) TGet(ctx context.Context, id uint, langID uint) (*Relation, error) {
	entity, err := s.repository.TGet(ctx, id, langID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *relationService) TFirst(ctx context.Context, entity *Relation, langID uint) (*Relation, error) {
	entity, err := s.repository.TFirst(ctx, entity, langID)
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

func (s *relationService) TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]Relation, error) {
	items, err := s.repository.TQuery(ctx, cond, langID)
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

func (s *relationService) UpsertBySysname(ctx context.Context, entity *Relation, langID uint) (err error) {

	found, err := s.propertyService.First(ctx, &property.Property{
		Sysname: entity.Sysname,
	})

	if err != nil {
		if err != yaruserror.ErrNotFound {
			return err
		}
		err = s.TCreate(ctx, entity, langID)
	} else {
		entity.ID = found.ID
		entity.NameSourceID = found.NameSourceID
		entity.DescriptionSourceID = found.DescriptionSourceID
		entity.CreatedAt = found.CreatedAt
		if found.SortOrder != property.SortOrderDefault {
			entity.SortOrder = found.SortOrder
		}
		err = s.TUpdate(ctx, entity, langID)
	}

	return err
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

func (s *relationService) TCreate(ctx context.Context, entity *Relation, langID uint) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.TCreate(ctx, entity, langID)
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

func (s *relationService) TUpdate(ctx context.Context, entity *Relation, langID uint) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.TUpdate(ctx, entity, langID)
	if err != nil {
		return errors.Wrapf(err, "Can not update an entity: %v", entity)
	}
	return nil
}

func (s *relationService) Delete(ctx context.Context, entity *Relation) error {
	err := s.repository.Delete(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not delete an entity: %v", entity)
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

func (s *relationService) TPropertyAndRelationQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]property.Property, []Relation, error) {
	props, rels, err := s.repository.TPropertyAndRelationQuery(ctx, cond, langID)
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

func (s *relationService) TGetPropertiesAndRelationsByEntityTypeID(ctx context.Context, entityTypeID uint, langID uint) ([]property.Property, []Relation, error) {
	props, rels, err := s.repository.TGetPropertiesAndRelationsByEntityTypeID(ctx, entityTypeID, langID)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Can not get properties and relations by entityTypeID: %v", entityTypeID)
	}
	return props, rels, nil
}
