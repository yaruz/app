package yaruzplatform

import (
	"context"
	"github.com/pkg/errors"

	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"

	"github.com/yaruz/app/internal/domain/sn_account"
)

// SnAccountRepository is a repository for the model entity
type SnAccountRepository struct {
	repository
}

var _ sn_account.Repository = (*SnAccountRepository)(nil)

// NewSnAccountRepository creates a new SnAccountRepository
func NewSnAccountRepository(repository *repository) (*SnAccountRepository, error) {
	return &SnAccountRepository{repository: *repository}, nil
}

func (r *SnAccountRepository) New(ctx context.Context) (*sn_account.SNAccount, error) {
	entity, err := r.repository.NewEntityByEntityType(ctx, sn_account.EntityType)
	if err != nil {
		return nil, err
	}
	return &sn_account.SNAccount{
		Entity: entity,
	}, nil
}

func (r *SnAccountRepository) instantiate(ctx context.Context, entity *entity.Entity) (*sn_account.SNAccount, error) {
	entity.PropertyFinder = r.GetPropertyFinder()
	obj := &sn_account.SNAccount{
		ID:     entity.ID,
		Entity: entity,
	}

	SNIDPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, sn_account.PropertySysnameSNID)
	if err != nil {
		return nil, err
	}
	SNIDVal, ok := obj.PropertiesValues[SNIDPropID]
	if ok {
		accountID, err := property.GetValueText(SNIDVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "SnAccountRepository.instantiate error. ")
		}
		obj.SNID = accountID
	}

	return obj, nil
}

// Get reads the album with the specified ID from the database.
func (r *SnAccountRepository) Get(ctx context.Context, id uint, langID uint) (*sn_account.SNAccount, error) {
	entityTypeID, err := r.yaruzRepository.ReferenceSubsystem().EntityType.GetIDBySysname(ctx, sn_account.EntityType)
	if err != nil {
		return nil, err
	}

	e, err := r.yaruzRepository.DataSubsystem().Entity.Get(ctx, id, entityTypeID, langID)
	if err != nil {
		return nil, err
	}

	return r.instantiate(ctx, e)
}

func (r *SnAccountRepository) First(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (*sn_account.SNAccount, error) {
	e, err := r.yaruzRepository.DataSubsystem().Entity.First(ctx, condition, &sn_account.SNAccount{}, langID)
	if err != nil {
		return nil, err
	}
	return r.instantiate(ctx, e)
}

// Query retrieves records with the specified offset and limit from the database.
func (r *SnAccountRepository) Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]sn_account.SNAccount, error) {

	entities, err := r.yaruzRepository.DataSubsystem().Entity.Query(ctx, condition, &sn_account.SNAccount{}, langID)
	if err != nil {
		return nil, err
	}

	items := make([]sn_account.SNAccount, len(entities))
	for i, e := range entities {
		obj, err := r.instantiate(ctx, &e)
		if err != nil {
			return nil, err
		}
		items[i] = *obj
	}

	return items, nil
}

func (r *SnAccountRepository) Count(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (uint, error) {
	return r.yaruzRepository.DataSubsystem().Entity.Count(ctx, condition, &sn_account.SNAccount{}, langID)
}

// Create saves a new record in the database.
func (r *SnAccountRepository) Create(ctx context.Context, obj *sn_account.SNAccount, langID uint) error {
	if obj.ID > 0 {
		return errors.New("entity is not new")
	}

	err := r.yaruzRepository.DataSubsystem().Entity.Create(ctx, obj.Entity, langID)
	obj.ID = obj.Entity.ID
	return err
}

// Update saves a changed Maintenance record in the database.
func (r *SnAccountRepository) Update(ctx context.Context, obj *sn_account.SNAccount, langID uint) error {
	if obj.ID == 0 {
		return errors.New("entity is new")
	}

	return r.yaruzRepository.DataSubsystem().Entity.Update(ctx, obj.Entity, langID)
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *SnAccountRepository) Delete(ctx context.Context, id uint) error {
	entityTypeID, err := r.yaruzRepository.ReferenceSubsystem().EntityType.GetIDBySysname(ctx, sn_account.EntityType)
	if err != nil {
		return err
	}

	return r.yaruzRepository.DataSubsystem().Entity.Delete(ctx, id, entityTypeID)
}
