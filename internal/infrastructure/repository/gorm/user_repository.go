package gorm

import (
	"context"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/yaruz/app/internal/pkg/apperror"

	"github.com/yaruz/app/internal/domain/user"
)

// UserRepository is a repository for the user entity
type UserRepository struct {
	repository
}

var _ user.Repository = (*UserRepository)(nil)

// New creates a new UserRepository
func NewUserRepository(repository *repository) (*UserRepository, error) {
	r := &UserRepository{repository: *repository}
	r.autoMigrate()
	return r, nil
}

func (r UserRepository) autoMigrate() {
	if r.db.IsAutoMigrate() {
		r.db.DB().AutoMigrate(&user.User{})
	}
}

// Get reads the album with the specified ID from the database.
func (r UserRepository) Get(ctx context.Context, id uint) (*user.User, error) {
	entity := &user.User{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

func (r UserRepository) First(ctx context.Context, entity *user.User) (*user.User, error) {
	err := r.DB().Where(entity).First(entity).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, apperror.ErrNotFound
		}
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r UserRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]user.User, error) {
	items := []user.User{}

	db := minipkg_gorm.Conditions(r.DB(), cond)
	if db.Error != nil {
		return nil, db.Error
	}

	err := db.Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return items, apperror.ErrNotFound
		}
	}
	return items, err
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r UserRepository) Create(ctx context.Context, entity *user.User) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}
