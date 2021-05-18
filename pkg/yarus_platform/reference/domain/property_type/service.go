package property_type

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type2property_view_type"

	"github.com/minipkg/selection_condition"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
)

// IService encapsulates usecase logic.
type IService interface {
	NewEntity() *PropertyType
	Get(ctx context.Context, id uint) (*PropertyType, error)
	Query(ctx context.Context, query *selection_condition.SelectionCondition) ([]PropertyType, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *PropertyType) error
	Update(ctx context.Context, entity *PropertyType) error
	Save(ctx context.Context, entity *PropertyType) error
	Delete(ctx context.Context, id uint) error
	InitPropertyViewTypes(ctx context.Context, entity *PropertyType) error
	BindPropertyViewType(ctx context.Context, entity *PropertyType, viewTypeID uint) error
	UnbindPropertyViewType(ctx context.Context, entity *PropertyType, viewTypeID uint) error
}

type service struct {
	logger                                  log.ILogger
	repository                              Repository
	propertyType2propertyViewTypeRepository property_type2property_view_type.Repository
}

var _ IService = (*service)(nil)

// NewService creates a new service.
func NewService(logger log.ILogger, repo Repository, propertyType2propertyViewTypeRepository property_type2property_view_type.Repository) IService {
	s := &service{
		logger:                                  logger,
		repository:                              repo,
		propertyType2propertyViewTypeRepository: propertyType2propertyViewTypeRepository,
	}
	repo.SetDefaultConditions(s.defaultConditions())
	return s
}

// Defaults returns defaults params
func (s *service) defaultConditions() *selection_condition.SelectionCondition {
	return &selection_condition.SelectionCondition{}
}

func (s *service) NewEntity() *PropertyType {
	return New()
}

// Get returns the entity with the specified ID.
func (s *service) Get(ctx context.Context, id uint) (*PropertyType, error) {
	entity, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Query returns the items with the specified selection condition.
func (s *service) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]PropertyType, error) {
	items, err := s.repository.Query(ctx, cond)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of items by query: %v", cond)
	}
	return items, nil
}

func (s *service) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	count, err := s.repository.Count(ctx, cond)
	if err != nil {
		return 0, errors.Wrapf(err, "Can not count a list of items by query: %v", cond)
	}
	return count, nil
}

func (s *service) Create(ctx context.Context, entity *PropertyType) error {
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

func (s *service) Update(ctx context.Context, entity *PropertyType) error {
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

func (s *service) Save(ctx context.Context, entity *PropertyType) error {
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

func (s *service) Delete(ctx context.Context, id uint) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "Can not delete an entity by ID: %v", id)
	}
	return nil
}

func (s *service) InitPropertyViewTypes(ctx context.Context, entity *PropertyType) error {
	return s.repository.InitPropertyViewTypes(ctx, entity)
}

func (s *service) BindPropertyViewType(ctx context.Context, entity *PropertyType, viewTypeID uint) error {
	return s.repository.BindPropertyViewType(ctx, entity, viewTypeID)
}

func (s *service) UnbindPropertyViewType(ctx context.Context, entity *PropertyType, viewTypeID uint) error {
	return s.repository.UnbindPropertyViewType(ctx, entity, viewTypeID)
}
