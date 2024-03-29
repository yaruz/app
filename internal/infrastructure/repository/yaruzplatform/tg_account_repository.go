package yaruzplatform

import (
	"context"
	"github.com/pkg/errors"
	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
	"time"

	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"

	"github.com/yaruz/app/internal/domain/tg_account"
)

// TgAccountRepository is a repository for the model entity
type TgAccountRepository struct {
	repository
}

var _ tg_account.Repository = (*TgAccountRepository)(nil)

// NewTgAccountRepository creates a new TgAccountRepository
func NewTgAccountRepository(repository *repository) (*TgAccountRepository, error) {
	return &TgAccountRepository{repository: *repository}, nil
}

func (r *TgAccountRepository) New(ctx context.Context) (*tg_account.TgAccount, error) {
	entity, err := r.repository.NewEntityByEntityType(ctx, tg_account.EntityType)
	if err != nil {
		return nil, err
	}
	return &tg_account.TgAccount{
		Entity: entity,
	}, nil
}

func (r *TgAccountRepository) instantiate(ctx context.Context, entity *entity.Entity) (*tg_account.TgAccount, error) {
	entity.PropertyFinder = r.GetPropertyFinder()
	obj := &tg_account.TgAccount{
		ID:     entity.ID,
		Entity: entity,
	}

	userIDPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, tg_account.PropertySysname_UserID)
	if err != nil {
		return nil, err
	}
	userIDVal, ok := obj.PropertiesValues[userIDPropID]
	if ok {
		userID, err := property.GetValueInt(userIDVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "TgAccountRepository.instantiate error. ")
		}
		obj.UserID = userID
	}

	accessHashPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, tg_account.PropertySysname_AccessHash)
	if err != nil {
		return nil, err
	}
	accessHashVal, ok := obj.PropertiesValues[accessHashPropID]
	if ok {
		accessHash, err := property.GetValueInt(accessHashVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "TgAccountRepository.instantiate error. ")
		}
		obj.AccessHash = accessHash
	}

	firstNamePropID, err := obj.PropertyFinder.GetIDBySysname(ctx, tg_account.PropertySysname_FirstName)
	if err != nil {
		return nil, err
	}
	firstNameVal, ok := obj.PropertiesValues[firstNamePropID]
	if ok {
		firstName, err := property.GetValueText(firstNameVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "TgAccountRepository.instantiate error. ")
		}
		obj.FirstName = firstName
	}

	lastNamePropID, err := obj.PropertyFinder.GetIDBySysname(ctx, tg_account.PropertySysname_LastName)
	if err != nil {
		return nil, err
	}
	lastNameVal, ok := obj.PropertiesValues[lastNamePropID]
	if ok {
		lastName, err := property.GetValueText(lastNameVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "TgAccountRepository.instantiate error. ")
		}
		obj.LastName = lastName
	}

	userNamePropID, err := obj.PropertyFinder.GetIDBySysname(ctx, tg_account.PropertySysname_UserName)
	if err != nil {
		return nil, err
	}
	userNameVal, ok := obj.PropertiesValues[userNamePropID]
	if ok {
		userName, err := property.GetValueText(userNameVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "TgAccountRepository.instantiate error. ")
		}
		obj.UserName = userName
	}

	phonePropID, err := obj.PropertyFinder.GetIDBySysname(ctx, tg_account.PropertySysname_Phone)
	if err != nil {
		return nil, err
	}
	phoneVal, ok := obj.PropertiesValues[phonePropID]
	if ok {
		phone, err := property.GetValueText(phoneVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "TgAccountRepository.instantiate error. ")
		}
		obj.Phone = phone
	}

	photoPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, tg_account.PropertySysname_Photo)
	if err != nil {
		return nil, err
	}
	photoVal, ok := obj.PropertiesValues[photoPropID]
	if ok {
		photo, err := property.GetValueText(photoVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "TgAccountRepository.instantiate error. ")
		}
		obj.Photo = photo
	}

	langCodePropID, err := obj.PropertyFinder.GetIDBySysname(ctx, tg_account.PropertySysname_LangCode)
	if err != nil {
		return nil, err
	}
	langCodeVal, ok := obj.PropertiesValues[langCodePropID]
	if ok {
		langCode, err := property.GetValueText(langCodeVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "TgAccountRepository.instantiate error. ")
		}
		obj.LangCode = langCode
	}

	createdAtPropID, err := obj.PropertyFinder.GetIDBySysname(ctx, tg_account.PropertySysname_CreatedAt)
	if err != nil {
		return nil, err
	}
	createdAtVal, ok := obj.PropertiesValues[createdAtPropID]
	if ok {
		createdAt, err := property.GetValueTime(createdAtVal.Value)
		if err != nil {
			return nil, errors.Wrapf(err, "TgAccountRepository.instantiate error. ")
		}
		obj.CreatedAt = createdAt
	}

	return obj, nil
}

// Get reads the album with the specified ID from the database.
func (r *TgAccountRepository) Get(ctx context.Context, id uint, langID uint) (*tg_account.TgAccount, error) {
	entityTypeID, err := r.yaruzRepository.ReferenceSubsystem().EntityType.GetIDBySysname(ctx, tg_account.EntityType)
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

func (r *TgAccountRepository) First(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (*tg_account.TgAccount, error) {
	e, err := r.yaruzRepository.DataSubsystem().Entity.First(ctx, condition, &tg_account.TgAccount{}, langID)
	if err != nil {
		if errors.Is(err, yaruserror.ErrNotFound) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return r.instantiate(ctx, e)
}

// Query retrieves records with the specified offset and limit from the database.
func (r *TgAccountRepository) Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]tg_account.TgAccount, error) {

	entities, err := r.yaruzRepository.DataSubsystem().Entity.Query(ctx, condition, &tg_account.TgAccount{}, langID)
	if err != nil {
		if errors.Is(err, yaruserror.ErrNotFound) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	items := make([]tg_account.TgAccount, len(entities))
	for i, e := range entities {
		obj, err := r.instantiate(ctx, &e)
		if err != nil {
			return nil, err
		}
		items[i] = *obj
	}

	return items, nil
}

func (r *TgAccountRepository) Count(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (uint, error) {
	return r.yaruzRepository.DataSubsystem().Entity.Count(ctx, condition, &tg_account.TgAccount{}, langID)
}

// Create saves a new record in the database.
func (r *TgAccountRepository) Create(ctx context.Context, obj *tg_account.TgAccount, langID uint) error {
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
func (r *TgAccountRepository) Update(ctx context.Context, obj *tg_account.TgAccount, langID uint) error {
	if obj.ID == 0 {
		return errors.New("entity is new")
	}

	return r.yaruzRepository.DataSubsystem().Entity.Update(ctx, obj.Entity, langID)
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *TgAccountRepository) Delete(ctx context.Context, id uint) error {
	entityTypeID, err := r.yaruzRepository.ReferenceSubsystem().EntityType.GetIDBySysname(ctx, tg_account.EntityType)
	if err != nil {
		if errors.Is(err, yaruserror.ErrNotFound) {
			return apperror.ErrNotFound
		}
		return err
	}

	return r.yaruzRepository.DataSubsystem().Entity.Delete(ctx, id, entityTypeID)
}
