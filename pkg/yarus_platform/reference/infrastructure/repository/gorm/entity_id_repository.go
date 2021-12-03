package gorm

import "golang.org/x/tools/go/ssa/interp/testdata/src/fmt"

type EntityIDRepository struct {
	repository
	entityTypes []string
}

func NewEntityIDRepository(repository *repository) *EntityIDRepository {
	return &EntityIDRepository{
		repository: *repository,
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

func (r *EntityIDRepository) getSeqName(typeSysname string) string {
	return typeSysname
}

func (r *EntityIDRepository) createQuery(name string) string {
	return fmt.Sprintf("CREATE SEQUENCE IF NOT EXISTS %s", name)
}
