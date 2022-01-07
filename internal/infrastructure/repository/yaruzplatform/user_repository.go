package yaruzplatform

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"

	"github.com/yaruz/app/internal/domain/user"

	"github.com/minipkg/selection_condition"
)

// UserRepository is a repository for the model entity
type UserRepository struct {
	repository
}

var _ user.Repository = (*UserRepository)(nil)

// New creates a new UserRepository
func NewUserRepository(repository *repository) (*UserRepository, error) {
	return &UserRepository{repository: *repository}, nil
}

func (r *UserRepository) New(ctx context.Context) (*user.User, error) {
	entity, err := r.repository.NewEntityByEntityType(ctx, user.EntityType)
	if err != nil {
		return nil, err
	}
	return &user.User{
		Entity: entity,
	}, nil
}

func (r *UserRepository) instantiate(ctx context.Context, entity *entity.Entity) (*user.User, error) {
	entity.PropertyFinder = r.GetPropertyFinder()
	obj := &user.User{
		Entity: entity,
	}

	namePropID, err := obj.PropertyFinder.GetIDBySysname(ctx, user.PropertySysnameEmail)
	if err != nil {
		return nil, err
	}
	nameVal, ok := obj.PropertiesValues[namePropID]
	if ok {
		if obj.Name, err = property.GetValueText(nameVal.Value); err != nil {
			return nil, errors.Wrapf(err, "UserRepository.instantiate error. ")
		}
	}

	agePropID, err := obj.PropertyFinder.GetIDBySysname(ctx, user.PropertySysnamePhone)
	if err != nil {
		return nil, err
	}
	ageVal, ok := obj.PropertiesValues[agePropID]
	if ok {
		age, err := property.GetValueInt(ageVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "UserRepository.instantiate error. ")
		}
		obj.Age = uint(age)
	}

	heightPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, user.PropertySysnameHeight)
	if err != nil {
		return nil, err
	}
	heightVal, ok := obj.PropertiesValues[heightPropID]
	if ok {
		if obj.Height, err = property.GetValueFloat(heightVal.Value); err != nil {
			return nil, errors.Wrapf(err, "UserRepository.instantiate error. ")
		}
	}

	weightPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, user.PropertySysnameWeight)
	if err != nil {
		return nil, err
	}
	weightVal, ok := obj.PropertiesValues[weightPropID]
	if ok {
		if obj.Weight, err = property.GetValueFloat(weightVal.Value); err != nil {
			return nil, errors.Wrapf(err, "UserRepository.instantiate error. ")
		}
	}

	return obj, nil
}

// Get reads the album with the specified ID from the database.
func (r *UserRepository) Get(ctx context.Context, id uint, langID uint) (*user.User, error) {
	entityTypeID, err := r.yaruzRepository.ReferenceSubsystem().EntityType.GetIDBySysname(ctx, user.EntityType)
	if err != nil {
		return nil, err
	}

	e, err := r.yaruzRepository.DataSubsystem().Entity.Get(ctx, id, entityTypeID, langID)
	if err != nil {
		return nil, err
	}

	return r.instantiate(ctx, e)
}

func (r *UserRepository) First(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (*user.User, error) {
	e, err := r.yaruzRepository.DataSubsystem().Entity.First(ctx, condition, &user.User{}, langID)
	if err != nil {
		return nil, err
	}
	return r.instantiate(ctx, e)
}

// Query retrieves records with the specified offset and limit from the database.
func (r *UserRepository) Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]user.User, error) {

	entities, err := r.yaruzRepository.DataSubsystem().Entity.Query(ctx, condition, &user.User{}, langID)
	if err != nil {
		return nil, err
	}

	items := make([]user.User, len(entities))
	for i, e := range entities {
		obj, err := r.instantiate(ctx, &e)
		if err != nil {
			return nil, err
		}
		items[i] = *obj
	}

	return items, nil
}

func (r *UserRepository) Count(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (uint, error) {
	return r.yaruzRepository.DataSubsystem().Entity.Count(ctx, condition, &user.User{}, langID)
}

// Create saves a new record in the database.
func (r *UserRepository) Create(ctx context.Context, obj *user.User, langID uint) error {
	if obj.ID > 0 {
		return errors.New("entity is not new")
	}

	return r.yaruzRepository.DataSubsystem().Entity.Create(ctx, obj.Entity, langID)
}

// Update saves a changed Maintenance record in the database.
func (r *UserRepository) Update(ctx context.Context, obj *user.User, langID uint) error {
	if obj.ID == 0 {
		return errors.New("entity is new")
	}

	return r.yaruzRepository.DataSubsystem().Entity.Update(ctx, obj.Entity, langID)
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	entityTypeID, err := r.yaruzRepository.ReferenceSubsystem().EntityType.GetIDBySysname(ctx, user.EntityType)
	if err != nil {
		return err
	}

	return r.yaruzRepository.DataSubsystem().Entity.Delete(ctx, id, entityTypeID)
}
