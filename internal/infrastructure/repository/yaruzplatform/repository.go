package yaruzplatform

import (
	"context"
	"github.com/yaruz/app/internal/domain/advertiser"
	"github.com/yaruz/app/internal/domain/advertising_campaign"
	"github.com/yaruz/app/internal/domain/offer"
	"github.com/yaruz/app/internal/domain/tg_account"

	"github.com/minipkg/selection_condition"
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
	case tg_account.EntityType:
		repo, err = NewTgAccountRepository(r)
	case advertiser.EntityType:
		repo, err = NewAdvertiserRepository(r)
	case advertising_campaign.EntityType:
		repo, err = NewAdvertisingCampaignRepository(r)
	case offer.EntityType:
		repo, err = NewOfferRepository(r)
	default:
		err = errors.Errorf("Case for entity %q not found", entity)
	}
	return repo, err
}

func (r *repository) SetDefaultConditions(defaultConditions *selection_condition.SelectionCondition) {
	r.Conditions = defaultConditions
}

func (r *repository) GetPropertyFinder() entity.PropertyFinder {
	return r.yaruzRepository.ReferenceSubsystem().Property
}

func (r *repository) NewEntityByEntityType(ctx context.Context, entityType string) (*entity.Entity, error) {
	entityTypeID, err := r.yaruzRepository.ReferenceSubsystem().EntityType.GetIDBySysname(ctx, entityType)
	if err != nil {
		return nil, err
	}

	entity := entity.New()
	entity.EntityTypeID = entityTypeID
	entity.PropertyFinder = r.GetPropertyFinder()

	return entity, nil
}
