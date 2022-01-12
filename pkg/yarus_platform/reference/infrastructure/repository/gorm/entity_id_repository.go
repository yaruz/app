package gorm

import (
	"fmt"

	"github.com/yaruz/app/pkg/yarus_platform/config"

	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
	"gorm.io/gorm"
)

type EntityIDRepository struct {
	repository
	entityTypesByClusterSysnames map[string][]string
	clusterSysnamesByEntityTypes map[string]string
}

func NewEntityIDRepository(repository *repository, entityTypesByClusterSysnames map[string][]string) *EntityIDRepository {
	clusterSysnamesByEntityTypes := make(map[string]string)

	for sysname, entityTypes := range entityTypesByClusterSysnames {
		for _, entityType := range entityTypes {
			clusterSysnamesByEntityTypes[entityType] = sysname
		}
	}
	entityTypesByClusterSysnames[config.DBClusterDefaultSysname] = nil

	return &EntityIDRepository{
		repository:                   *repository,
		entityTypesByClusterSysnames: entityTypesByClusterSysnames,
		clusterSysnamesByEntityTypes: clusterSysnamesByEntityTypes,
	}
}

func (r *EntityIDRepository) AutoMigrate() error {
	for s := range r.entityTypesByClusterSysnames {
		query := r.createQuery(r.getSeqNameByClusterSysname(s))
		if err := r.db.DB().Exec(query).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *EntityIDRepository) isSeparatedEntityType(entityTypeSysname string) (string, bool) {
	sysname, ok := r.clusterSysnamesByEntityTypes[entityTypeSysname]
	return sysname, ok
}

func (r *EntityIDRepository) getSeqNameByClusterSysname(clusterSysname string) string {
	return r.buildSeqName(clusterSysname)
}

func (r *EntityIDRepository) getSeqNameByEntityTypeSysname(entityTypeSysname string) string {
	if clusterSysname, ok := r.isSeparatedEntityType(entityTypeSysname); ok {
		return r.buildSeqName(clusterSysname)
	}
	return r.buildSeqName(config.DBClusterDefaultSysname)
}

func (r *EntityIDRepository) buildSeqName(clusterSysname string) string {
	return "id4type_" + clusterSysname
}

func (r *EntityIDRepository) NextVal(entityTypeSysname string) (id uint, err error) {
	return r.getVal(r.nextvalQuery(r.getSeqNameByEntityTypeSysname(entityTypeSysname)))
}

func (r *EntityIDRepository) LastVal(entityTypeSysname string) (id uint, err error) {
	return r.getVal(r.lastvalQuery(r.getSeqNameByEntityTypeSysname(entityTypeSysname)))
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
	return fmt.Sprintf("CREATE SEQUENCE IF NOT EXISTS %s", name)
}

func (r *EntityIDRepository) nextvalQuery(name string) string {
	return fmt.Sprintf("SELECT nextval('%s')", name)
}

func (r *EntityIDRepository) lastvalQuery(name string) string {
	return fmt.Sprintf("SELECT lastval('%s')", name)
}
