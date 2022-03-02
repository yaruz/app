package yaruzplatform

import (
	"context"
	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"

	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"

	"github.com/yaruz/app/internal/domain/offer"

	"github.com/minipkg/selection_condition"
)

// OfferRepository is a repository for the model entity
type OfferRepository struct {
	repository
}

var _ offer.Repository = (*OfferRepository)(nil)

// New creates a new OfferRepository
func NewOfferRepository(repository *repository) (*OfferRepository, error) {
	return &OfferRepository{repository: *repository}, nil
}

func (r *OfferRepository) New(ctx context.Context) (*offer.Offer, error) {
	entity, err := r.repository.NewEntityByEntityType(ctx, offer.EntityType)
	if err != nil {
		return nil, err
	}
	return &offer.Offer{
		Entity: entity,
	}, nil
}

func (r *OfferRepository) instantiate(ctx context.Context, entity *entity.Entity) (*offer.Offer, error) {
	entity.PropertyFinder = r.GetPropertyFinder()
	obj := &offer.Offer{
		ID:     entity.ID,
		Entity: entity,
	}

	createdAtPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, offer.PropertySysnameCreatedAt)
	if err != nil {
		return nil, err
	}
	createdAtVal, ok := obj.PropertiesValues[createdAtPropID]
	if ok {
		if obj.CreatedAt, err = property.GetValueTime(createdAtVal.Value); err != nil {
			return nil, errors.Wrapf(err, "OfferRepository.instantiate error. ")
		}
	}

	startedAtPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, offer.PropertySysnameStartedAt)
	if err != nil {
		return nil, err
	}
	startedAtVal, ok := obj.PropertiesValues[startedAtPropID]
	if ok {
		if obj.StartedAt, err = property.GetValueTime(startedAtVal.Value); err != nil {
			return nil, errors.Wrapf(err, "OfferRepository.instantiate error. ")
		}
	}

	finishedAtPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, offer.PropertySysnameFinishedAt)
	if err != nil {
		return nil, err
	}
	finishedAtVal, ok := obj.PropertiesValues[finishedAtPropID]
	if ok {
		if obj.FinishedAt, err = property.GetValueTime(finishedAtVal.Value); err != nil {
			return nil, errors.Wrapf(err, "OfferRepository.instantiate error. ")
		}
	}

	return obj, nil
}

// Get reads the album with the specified ID from the database.
func (r *OfferRepository) Get(ctx context.Context, id uint, langID uint) (*offer.Offer, error) {
	entityTypeID, err := r.yaruzRepository.ReferenceSubsystem().EntityType.GetIDBySysname(ctx, offer.EntityType)
	if err != nil {
		if errors.Is(err, yaruserror.ErrNotFound) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	e, err := r.yaruzRepository.DataSubsystem().Entity.Get(ctx, id, entityTypeID, langID)
	if err != nil {
		if errors.Is(err, yaruserror.ErrNotFound) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	return r.instantiate(ctx, e)
}

func (r *OfferRepository) First(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (*offer.Offer, error) {
	e, err := r.yaruzRepository.DataSubsystem().Entity.First(ctx, condition, &offer.Offer{}, langID)
	if err != nil {
		if errors.Is(err, yaruserror.ErrNotFound) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return r.instantiate(ctx, e)
}

// Query retrieves records with the specified offset and limit from the database.
func (r *OfferRepository) Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]offer.Offer, error) {

	entities, err := r.yaruzRepository.DataSubsystem().Entity.Query(ctx, condition, &offer.Offer{}, langID)
	if err != nil {
		if errors.Is(err, yaruserror.ErrNotFound) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	items := make([]offer.Offer, len(entities))
	for i, e := range entities {
		obj, err := r.instantiate(ctx, &e)
		if err != nil {
			return nil, err
		}
		items[i] = *obj
	}

	return items, nil
}

func (r *OfferRepository) Count(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (uint, error) {
	return r.yaruzRepository.DataSubsystem().Entity.Count(ctx, condition, &offer.Offer{}, langID)
}

// Create saves a new record in the database.
func (r *OfferRepository) Create(ctx context.Context, obj *offer.Offer, langID uint) error {
	if obj.ID > 0 {
		return errors.New("entity is not new")
	}

	err := r.yaruzRepository.DataSubsystem().Entity.Create(ctx, obj.Entity, langID)
	obj.ID = obj.Entity.ID
	return err
}

// Update saves a changed Maintenance record in the database.
func (r *OfferRepository) Update(ctx context.Context, obj *offer.Offer, langID uint) error {
	if obj.ID == 0 {
		return errors.New("entity is new")
	}

	return r.yaruzRepository.DataSubsystem().Entity.Update(ctx, obj.Entity, langID)
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *OfferRepository) Delete(ctx context.Context, id uint) error {
	entityTypeID, err := r.yaruzRepository.ReferenceSubsystem().EntityType.GetIDBySysname(ctx, offer.EntityType)
	if err != nil {
		if errors.Is(err, yaruserror.ErrNotFound) {
			return apperror.ErrNotFound
		}
		return err
	}

	return r.yaruzRepository.DataSubsystem().Entity.Delete(ctx, id, entityTypeID)
}
