package yaruzplatform

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"

	"github.com/yaruz/app/internal/domain/advertiser"

	"github.com/minipkg/selection_condition"
)

// AdvertiserRepository is a repository for the model entity
type AdvertiserRepository struct {
	repository
}

var _ advertiser.Repository = (*AdvertiserRepository)(nil)

// New creates a new AdvertiserRepository
func NewAdvertiserRepository(repository *repository) (*AdvertiserRepository, error) {
	return &AdvertiserRepository{repository: *repository}, nil
}

func (r *AdvertiserRepository) New(ctx context.Context) (*advertiser.Advertiser, error) {
	entity, err := r.repository.NewEntityByEntityType(ctx, advertiser.EntityType)
	if err != nil {
		return nil, err
	}
	return &advertiser.Advertiser{
		Entity: entity,
	}, nil
}

func (r *AdvertiserRepository) instantiate(ctx context.Context, entity *entity.Entity) (*advertiser.Advertiser, error) {
	entity.PropertyFinder = r.GetPropertyFinder()
	obj := &advertiser.Advertiser{
		Entity: entity,
	}

	namePropID, err := obj.PropertyFinder.GetIDBySysname(ctx, advertiser.PropertySysnameEmail)
	if err != nil {
		return nil, err
	}
	nameVal, ok := obj.PropertiesValues[namePropID]
	if ok {
		if obj.Name, err = property.GetValueText(nameVal.Value); err != nil {
			return nil, errors.Wrapf(err, "AdvertiserRepository.instantiate error. ")
		}
	}

	agePropID, err := obj.PropertyFinder.GetIDBySysname(ctx, advertiser.PropertySysnamePhone)
	if err != nil {
		return nil, err
	}
	ageVal, ok := obj.PropertiesValues[agePropID]
	if ok {
		age, err := property.GetValueInt(ageVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "AdvertiserRepository.instantiate error. ")
		}
		obj.Age = uint(age)
	}

	heightPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, advertiser.PropertySysnameHeight)
	if err != nil {
		return nil, err
	}
	heightVal, ok := obj.PropertiesValues[heightPropID]
	if ok {
		if obj.Height, err = property.GetValueFloat(heightVal.Value); err != nil {
			return nil, errors.Wrapf(err, "AdvertiserRepository.instantiate error. ")
		}
	}

	weightPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, advertiser.PropertySysnameWeight)
	if err != nil {
		return nil, err
	}
	weightVal, ok := obj.PropertiesValues[weightPropID]
	if ok {
		if obj.Weight, err = property.GetValueFloat(weightVal.Value); err != nil {
			return nil, errors.Wrapf(err, "AdvertiserRepository.instantiate error. ")
		}
	}

	return obj, nil
}

// Get reads the album with the specified ID from the database.
func (r *AdvertiserRepository) Get(ctx context.Context, id uint, langID uint) (*advertiser.Advertiser, error) {
	entityTypeID, err := r.yaruzRepository.ReferenceSubsystem().EntityType.GetIDBySysname(ctx, advertiser.EntityType)
	if err != nil {
		return nil, err
	}

	e, err := r.yaruzRepository.DataSubsystem().Entity.Get(ctx, id, entityTypeID, langID)
	if err != nil {
		return nil, err
	}

	return r.instantiate(ctx, e)
}

func (r *AdvertiserRepository) First(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (*advertiser.Advertiser, error) {
	e, err := r.yaruzRepository.DataSubsystem().Entity.First(ctx, condition, &advertiser.Advertiser{}, langID)
	if err != nil {
		return nil, err
	}
	return r.instantiate(ctx, e)
}

// Query retrieves records with the specified offset and limit from the database.
func (r *AdvertiserRepository) Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]advertiser.Advertiser, error) {

	entities, err := r.yaruzRepository.DataSubsystem().Entity.Query(ctx, condition, &advertiser.Advertiser{}, langID)
	if err != nil {
		return nil, err
	}

	items := make([]advertiser.Advertiser, len(entities))
	for i, e := range entities {
		obj, err := r.instantiate(ctx, &e)
		if err != nil {
			return nil, err
		}
		items[i] = *obj
	}

	return items, nil
}

func (r *AdvertiserRepository) Count(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (uint, error) {
	return r.yaruzRepository.DataSubsystem().Entity.Count(ctx, condition, &advertiser.Advertiser{}, langID)
}

// Create saves a new record in the database.
func (r *AdvertiserRepository) Create(ctx context.Context, obj *advertiser.Advertiser, langID uint) error {
	if obj.ID > 0 {
		return errors.New("entity is not new")
	}

	return r.yaruzRepository.DataSubsystem().Entity.Create(ctx, obj.Entity, langID)
}

// Update saves a changed Maintenance record in the database.
func (r *AdvertiserRepository) Update(ctx context.Context, obj *advertiser.Advertiser, langID uint) error {
	if obj.ID == 0 {
		return errors.New("entity is new")
	}

	return r.yaruzRepository.DataSubsystem().Entity.Update(ctx, obj.Entity, langID)
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *AdvertiserRepository) Delete(ctx context.Context, id uint) error {
	entityTypeID, err := r.yaruzRepository.ReferenceSubsystem().EntityType.GetIDBySysname(ctx, advertiser.EntityType)
	if err != nil {
		return err
	}

	return r.yaruzRepository.DataSubsystem().Entity.Delete(ctx, id, entityTypeID)
}
