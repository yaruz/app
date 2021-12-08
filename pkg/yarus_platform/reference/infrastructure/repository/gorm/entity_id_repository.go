package gorm

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
	"gorm.io/gorm"
)

const ShardsNameDefault = "default"

type EntityIDRepository struct {
	repository
	entityTypes   []string
	entityTypesSt map[string]struct{}
}

func NewEntityIDRepository(repository *repository, entityTypes []string) *EntityIDRepository {
	entityTypes = append(entityTypes, ShardsNameDefault)
	entityTypesSt := make(map[string]struct{}, len(entityTypes))
	for _, entityType := range entityTypes {
		entityTypesSt[entityType] = struct{}{}
	}

	return &EntityIDRepository{
		repository:    *repository,
		entityTypes:   entityTypes,
		entityTypesSt: entityTypesSt,
	}
}

func (r *EntityIDRepository) AutoMigrate() error {
	for _, t := range r.entityTypes {
		query := r.createQuery(r.getSeqName(t))
		if err := r.db.DB().Exec(query).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *EntityIDRepository) isSpecialEntityType(entityTypeSysname string) bool {
	_, ok := r.entityTypesSt[entityTypeSysname]
	return ok
}

func (r *EntityIDRepository) getSeqName(entityTypeSysname string) string {
	if r.isSpecialEntityType(entityTypeSysname) {
		return entityTypeSysname
	}
	return ShardsNameDefault
}

func (r *EntityIDRepository) NextVal(entityTypeSysname string) (id uint, err error) {
	return r.getVal(r.nextvalQuery(r.getSeqName(entityTypeSysname)))
}

func (r *EntityIDRepository) LastVal(entityTypeSysname string) (id uint, err error) {
	return r.getVal(r.lastvalQuery(r.getSeqName(entityTypeSysname)))
}

func (r *EntityIDRepository) getVal(sql string) (id uint, err error) {
	if err = r.db.DB().Raw(sql).Scan(&id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, yaruserror.ErrNotFound
		}
		return 0, err
	}
	return id, nil
}

func (r *EntityIDRepository) createQuery(name string) string {
	return fmt.Sprintf("CREATE SEQUENCE IF NOT EXISTS '%s'", name)
}

func (r *EntityIDRepository) nextvalQuery(name string) string {
	return fmt.Sprintf("SELECT nextval('%s')", name)
}

func (r *EntityIDRepository) lastvalQuery(name string) string {
	return fmt.Sprintf("SELECT lastval('%s')", name)
}
