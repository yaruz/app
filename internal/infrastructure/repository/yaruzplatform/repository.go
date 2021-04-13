package yaruzplatform

import (
	"github.com/minipkg/selection_condition"
	"github.com/yaruz/app/internal/domain/task"

	"github.com/yaruz/app/pkg/yarus_platform"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
)

// IRepository is an interface of repository
type IRepository interface{}

// repository persists albums in database
type repository struct {
	yaruzRepository yarus_platform.IRepository
	logger          log.ILogger
	Conditions      *selection_condition.SelectionCondition
}

const DefaultLimit = 1000

// GetRepository return a repository
func GetRepository(logger log.ILogger, yaruzRepository yarus_platform.IRepository, entity string) (repo IRepository, err error) {
	r := &repository{
		yaruzRepository: yaruzRepository,
		logger:          logger,
	}

	switch entity {
	case task.EntityName:
		repo, err = NewTaskRepository(r)
	default:
		err = errors.Errorf("Repository for entity %q not found", entity)
	}
	return repo, err
}

func (r *repository) SetDefaultConditions(defaultConditions *selection_condition.SelectionCondition) {
	r.Conditions = defaultConditions
}
