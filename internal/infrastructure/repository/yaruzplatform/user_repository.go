package yaruzplatform

import (
	"context"
	"github.com/pkg/errors"
	"time"

	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"

	"github.com/yaruz/app/internal/pkg/apperror"

	"github.com/yaruz/app/internal/domain/tg_account"
	"github.com/yaruz/app/internal/domain/user"
)

// UserRepository is a repository for the model entity
type UserRepository struct {
	repository
	tgAccountRepository tg_account.Repository
}

var _ user.Repository = (*UserRepository)(nil)

// NewUserRepository creates a new UserRepository
func NewUserRepository(repository *repository, tgAccountRepository tg_account.Repository) (*UserRepository, error) {
	return &UserRepository{
		repository:          *repository,
		tgAccountRepository: tgAccountRepository,
	}, nil
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
		ID:     entity.ID,
		Entity: entity,
	}

	emailPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, user.PropertySysnameEmail)
	if err != nil {
		return nil, err
	}
	emailVal, ok := obj.PropertiesValues[emailPropID]
	if ok {
		if obj.Email, err = property.GetValueText(emailVal.Value); err != nil {
			return nil, errors.Wrapf(err, "UserRepository.instantiate error. ")
		}
	}

	phonePropID, err := obj.PropertyFinder.GetIDBySysname(ctx, user.PropertySysnamePhone)
	if err != nil {
		return nil, err
	}
	phoneVal, ok := obj.PropertiesValues[phonePropID]
	if ok {
		if obj.Phone, err = property.GetValueText(phoneVal.Value); err != nil {
			return nil, errors.Wrapf(err, "UserRepository.instantiate error. ")
		}
	}

	accountIDPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, user.PropertySysnameAccountID)
	if err != nil {
		return nil, err
	}
	accountIDVal, ok := obj.PropertiesValues[accountIDPropID]
	if ok {
		accountID, err := property.GetValueText(accountIDVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "UserRepository.instantiate error. ")
		}
		obj.AccountID = accountID
	}

	createdAtPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, user.PropertySysnameCreatedAt)
	if err != nil {
		return nil, err
	}
	createdAtVal, ok := obj.PropertiesValues[createdAtPropID]
	if ok {
		createdAt, err := property.GetValueTime(createdAtVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "UserRepository.instantiate error. ")
		}
		obj.CreatedAt = createdAt
	}

	return obj, nil
}

// Get reads the album with the specified ID from the database.
func (r *UserRepository) Get(ctx context.Context, id uint, langID uint) (*user.User, error) {
	entityTypeID, err := r.yaruzRepository.ReferenceSubsystem().EntityType.GetIDBySysname(ctx, user.EntityType)
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

func (r *UserRepository) First(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (*user.User, error) {
	e, err := r.yaruzRepository.DataSubsystem().Entity.First(ctx, condition, &user.User{}, langID)
	if err != nil {
		if errors.Is(err, yaruserror.ErrNotFound) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return r.instantiate(ctx, e)
}

// Query retrieves records with the specified offset and limit from the database.
func (r *UserRepository) Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]user.User, error) {

	entities, err := r.yaruzRepository.DataSubsystem().Entity.Query(ctx, condition, &user.User{}, langID)
	if err != nil {
		if errors.Is(err, yaruserror.ErrNotFound) {
			return nil, apperror.ErrNotFound
		}
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

	if err := obj.SetCreatedAt(ctx, time.Now()); err != nil {
		return err
	}

	err := r.yaruzRepository.DataSubsystem().Entity.Create(ctx, obj.Entity, langID)
	obj.ID = obj.Entity.ID
	return err
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
		if errors.Is(err, yaruserror.ErrNotFound) {
			return apperror.ErrNotFound
		}
		return err
	}

	return r.yaruzRepository.DataSubsystem().Entity.Delete(ctx, id, entityTypeID)
}

func (r *UserRepository) LangIDValidate(ctx context.Context, langID uint) error {
	_, err := r.yaruzRepository.ReferenceSubsystem().TextLang.GetCfgnameByID(ctx, langID)
	return err
}

func (r *UserRepository) GetTgAccount(ctx context.Context, obj *user.User, langID uint) (*tg_account.TgAccount, error) {
	relID, err := r.GetPropertyFinder().GetIDBySysname(ctx, user.RelationSysnameTgAccount)
	if err != nil {
		return nil, err
	}

	relVal, ok := obj.RelationsValues[relID]
	if !ok || len(relVal.Value) == 0 {
		return nil, apperror.ErrNotFound
	}
	tgAccID := relVal.Value[0]

	return r.tgAccountRepository.Get(ctx, tgAccID, langID)
}
