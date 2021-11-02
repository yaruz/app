package yaruzplatform

import (
	"context"

	"github.com/minipkg/selection_condition"
	"github.com/yaruz/app/internal/domain/task"
	"github.com/yaruz/app/internal/domain/user"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"

	"github.com/yaruz/app/pkg/yarus_platform"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
)

// IRepository is an interface of repository
type IRepository interface{}

// repository persists albums in database
type repository struct {
	yaruzRepository yarus_platform.IPlatform
	logger          log.ILogger
	Conditions      *selection_condition.SelectionCondition
}

const DefaultLimit = 1000

// GetRepository return a repository
func GetRepository(logger log.ILogger, yaruzRepository yarus_platform.IPlatform, entity string) (repo IRepository, err error) {
	r := &repository{
		yaruzRepository: yaruzRepository,
		logger:          logger,
	}

	switch entity {
	case user.EntityType:
		repo, err = NewUserRepository(r)
	case task.EntityType:
		repo, err = NewTaskRepository(r)
	default:
		err = errors.Errorf("Text for entity %q not found", entity)
	}
	return repo, err
}

func (r *repository) SetDefaultConditions(defaultConditions *selection_condition.SelectionCondition) {
	r.Conditions = defaultConditions
}

func (r *repository) GetPropertyFinder() entity.PropertyFinder {
	return r.yaruzRepository.ReferenceSubsystem().Property
}

func (r *repository) NewByEntityTypeID(ctx context.Context, entityTypeID uint) (*user.User, error) {
	entity := entity.New()
	entity.EntityTypeID = entityTypeID
	entity.PropertyFinder = r.GetPropertyFinder()

	return &user.User{
		Entity: entity,
	}, nil
}
