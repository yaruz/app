package yaruzplatform

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"

	"github.com/yaruz/app/internal/domain/advertising_campaign"

	"github.com/minipkg/selection_condition"
)

// AdvertisingCampaignRepository is a repository for the model entity
type AdvertisingCampaignRepository struct {
	repository
}

var _ advertising_campaign.Repository = (*AdvertisingCampaignRepository)(nil)

// New creates a new AdvertisingCampaignRepository
func NewAdvertisingCampaignRepository(repository *repository) (*AdvertisingCampaignRepository, error) {
	return &AdvertisingCampaignRepository{repository: *repository}, nil
}

func (r *AdvertisingCampaignRepository) New(ctx context.Context) (*advertising_campaign.AdvertisingCampaign, error) {
	entity, err := r.repository.NewEntityByEntityType(ctx, advertising_campaign.EntityType)
	if err != nil {
		return nil, err
	}
	return &advertising_campaign.AdvertisingCampaign{
		Entity: entity,
	}, nil
}

func (r *AdvertisingCampaignRepository) instantiate(ctx context.Context, entity *entity.Entity) (*advertising_campaign.AdvertisingCampaign, error) {
	entity.PropertyFinder = r.GetPropertyFinder()
	obj := &advertising_campaign.AdvertisingCampaign{
		Entity: entity,
	}

	namePropID, err := obj.PropertyFinder.GetIDBySysname(ctx, advertising_campaign.PropertySysnameEmail)
	if err != nil {
		return nil, err
	}
	nameVal, ok := obj.PropertiesValues[namePropID]
	if ok {
		if obj.Name, err = property.GetValueText(nameVal.Value); err != nil {
			return nil, errors.Wrapf(err, "AdvertisingCampaignRepository.instantiate error. ")
		}
	}

	agePropID, err := obj.PropertyFinder.GetIDBySysname(ctx, advertising_campaign.PropertySysnamePhone)
	if err != nil {
		return nil, err
	}
	ageVal, ok := obj.PropertiesValues[agePropID]
	if ok {
		age, err := property.GetValueInt(ageVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "AdvertisingCampaignRepository.instantiate error. ")
		}
		obj.Age = uint(age)
	}

	heightPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, advertising_campaign.PropertySysnameHeight)
	if err != nil {
		return nil, err
	}
	heightVal, ok := obj.PropertiesValues[heightPropID]
	if ok {
		if obj.Height, err = property.GetValueFloat(heightVal.Value); err != nil {
			return nil, errors.Wrapf(err, "AdvertisingCampaignRepository.instantiate error. ")
		}
	}

	weightPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, advertising_campaign.PropertySysnameWeight)
	if err != nil {
		return nil, err
	}
	weightVal, ok := obj.PropertiesValues[weightPropID]
	if ok {
		if obj.Weight, err = property.GetValueFloat(weightVal.Value); err != nil {
			return nil, errors.Wrapf(err, "AdvertisingCampaignRepository.instantiate error. ")
		}
	}

	return obj, nil
}

// Get reads the album with the specified ID from the database.
func (r *AdvertisingCampaignRepository) Get(ctx context.Context, id uint, langID uint) (*advertising_campaign.AdvertisingCampaign, error) {
	entityTypeID, err := r.yaruzRepository.ReferenceSubsystem().EntityType.GetIDBySysname(ctx, advertising_campaign.EntityType)
	if err != nil {
		return nil, err
	}

	e, err := r.yaruzRepository.DataSubsystem().Entity.Get(ctx, id, entityTypeID, langID)
	if err != nil {
		return nil, err
	}

	return r.instantiate(ctx, e)
}

func (r *AdvertisingCampaignRepository) First(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (*advertising_campaign.AdvertisingCampaign, error) {
	e, err := r.yaruzRepository.DataSubsystem().Entity.First(ctx, condition, &advertising_campaign.AdvertisingCampaign{}, langID)
	if err != nil {
		return nil, err
	}
	return r.instantiate(ctx, e)
}

// Query retrieves records with the specified offset and limit from the database.
func (r *AdvertisingCampaignRepository) Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]advertising_campaign.AdvertisingCampaign, error) {

	entities, err := r.yaruzRepository.DataSubsystem().Entity.Query(ctx, condition, &advertising_campaign.AdvertisingCampaign{}, langID)
	if err != nil {
		return nil, err
	}

	items := make([]advertising_campaign.AdvertisingCampaign, len(entities))
	for i, e := range entities {
		obj, err := r.instantiate(ctx, &e)
		if err != nil {
			return nil, err
		}
		items[i] = *obj
	}

	return items, nil
}

func (r *AdvertisingCampaignRepository) Count(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (uint, error) {
	return r.yaruzRepository.DataSubsystem().Entity.Count(ctx, condition, &advertising_campaign.AdvertisingCampaign{}, langID)
}

// Create saves a new record in the database.
func (r *AdvertisingCampaignRepository) Create(ctx context.Context, obj *advertising_campaign.AdvertisingCampaign, langID uint) error {
	if obj.ID > 0 {
		return errors.New("entity is not new")
	}

	return r.yaruzRepository.DataSubsystem().Entity.Create(ctx, obj.Entity, langID)
}

// Update saves a changed Maintenance record in the database.
func (r *AdvertisingCampaignRepository) Update(ctx context.Context, obj *advertising_campaign.AdvertisingCampaign, langID uint) error {
	if obj.ID == 0 {
		return errors.New("entity is new")
	}

	return r.yaruzRepository.DataSubsystem().Entity.Update(ctx, obj.Entity, langID)
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *AdvertisingCampaignRepository) Delete(ctx context.Context, id uint) error {
	entityTypeID, err := r.yaruzRepository.ReferenceSubsystem().EntityType.GetIDBySysname(ctx, advertising_campaign.EntityType)
	if err != nil {
		return err
	}

	return r.yaruzRepository.DataSubsystem().Entity.Delete(ctx, id, entityTypeID)
}
