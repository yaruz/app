package gorm

import (
	"context"

	"gorm.io/gorm/clause"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/bool_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/date_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/float_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/int_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/time_value"

	"github.com/minipkg/selection_condition"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/yaruz/app/internal/pkg/apperror"
	domain_entity "github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"
)

// EntityRepository is a repository for the model entity
type EntityRepository struct {
	repository
	valueRepositories *domain_entity.ValueRepositories
}

var _ domain_entity.Repository = (*EntityRepository)(nil)

// New creates a new EntityRepository
func NewEntityRepository(repository *repository, valueRepositories *domain_entity.ValueRepositories) (*EntityRepository, error) {
	return &EntityRepository{
		repository:        *repository,
		valueRepositories: valueRepositories,
	}, nil
}

// Create saves a new record in the database.
func (r *EntityRepository) Create(ctx context.Context, entity *domain_entity.Entity, langID uint) error {
	db, err := r.mapReducer.GetDB(ctx, entity.EntityTypeID, entity.ID)
	if err != nil {
		return nil
	}

	return db.DB().Transaction(func(tx *gorm.DB) (err error) {

		if err := tx.Omit(clause.Associations).Create(entity).Error; err != nil {
			return err
		}

		return r.afterSaveTx(ctx, entity, langID, tx)
	})
}

// Update saves a changed Maintenance record in the database.
func (r *EntityRepository) Update(ctx context.Context, entity *domain_entity.Entity, langID uint) error {
	db, err := r.mapReducer.GetDB(ctx, entity.ID, entity.EntityTypeID)
	if err != nil {
		return err
	}
	return r.updateTx(ctx, entity, langID, db.DB())
}

func (r *EntityRepository) updateTx(ctx context.Context, entity *domain_entity.Entity, langID uint, tx *gorm.DB) error {
	if err := entity.Validate(); err != nil {
		return err
	}

	return tx.Transaction(func(tx *gorm.DB) (err error) {

		if err := tx.Omit(clause.Associations).Save(entity).Error; err != nil {
			return err
		}

		return r.afterSaveTx(ctx, entity, langID, tx)
	})
}

// Сохраняем все значения всех свойств из всех слайсов ...Values
func (r *EntityRepository) afterSaveTx(ctx context.Context, entity *domain_entity.Entity, langID uint, tx *gorm.DB) error {

	if err := r.valueRepositories.Bool.BatchSaveChangesTx(ctx, entity.ID, entity.BoolValues, langID, tx); err != nil {
		return err
	}

	if err := r.valueRepositories.Int.BatchSaveChangesTx(ctx, entity.ID, entity.IntValues, langID, tx); err != nil {
		return err
	}

	if err := r.valueRepositories.Float.BatchSaveChangesTx(ctx, entity.ID, entity.FloatValues, langID, tx); err != nil {
		return err
	}

	if err := r.valueRepositories.Date.BatchSaveChangesTx(ctx, entity.ID, entity.DateValues, langID, tx); err != nil {
		return err
	}

	if err := r.valueRepositories.Time.BatchSaveChangesTx(ctx, entity.ID, entity.TimeValues, langID, tx); err != nil {
		return err
	}

	if err := r.valueRepositories.Text.BatchSaveChangesTx(ctx, entity.ID, entity.TextValues, langID, tx); err != nil {
		return err
	}

	return nil
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *EntityRepository) Delete(ctx context.Context, ID uint, typeID uint) error {
	db, err := r.mapReducer.GetDB(ctx, ID, typeID)
	if err != nil {
		return err
	}

	return db.DB().Transaction(func(tx *gorm.DB) (err error) {

		if err = r.valueRepositories.Bool.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
			Where: &bool_value.BoolValue{
				EntityID: ID,
			},
		}, tx); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		if err = r.valueRepositories.Int.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
			Where: &int_value.IntValue{
				EntityID: ID,
			},
		}, tx); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		if err = r.valueRepositories.Float.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
			Where: &float_value.FloatValue{
				EntityID: ID,
			},
		}, tx); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		if err = r.valueRepositories.Date.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
			Where: &date_value.DateValue{
				EntityID: ID,
			},
		}, tx); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		if err = r.valueRepositories.Time.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
			Where: &time_value.TimeValue{
				EntityID: ID,
			},
		}, tx); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		if err = r.valueRepositories.Text.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
			Where: &text_value.TextValue{
				EntityID: ID,
			},
		}, tx); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		if err = tx.Delete(&domain_entity.Entity{}, ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.ErrNotFound
			}
		}
		return err
	})
}
