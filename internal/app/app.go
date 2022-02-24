package app

import (
	"context"
	"github.com/yaruz/app/internal/pkg/auth"
	golog "log"

	"github.com/pkg/errors"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/db/redis"
	"github.com/minipkg/db/redis/cache"
	"github.com/minipkg/log"

	"github.com/yaruz/app/pkg/yarus_platform"

	"github.com/yaruz/app/internal/domain/advertiser"
	"github.com/yaruz/app/internal/domain/advertising_campaign"
	"github.com/yaruz/app/internal/domain/offer"
	"github.com/yaruz/app/internal/domain/session"
	"github.com/yaruz/app/internal/domain/sn_account"
	"github.com/yaruz/app/internal/domain/user"
	redisrepo "github.com/yaruz/app/internal/infrastructure/repository/redis"
	"github.com/yaruz/app/internal/infrastructure/repository/yaruzplatform"
	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/internal/pkg/config"
)

// App struct is the common part of all applications
type App struct {
	Cfg    config.Configuration
	Domain Domain
	Auth   Auth
	Infra  *Infrastructure
}

type Infrastructure struct {
	Logger          log.ILogger
	IdentityDB      minipkg_gorm.IDB
	Redis           redis.IDB
	Cache           cache.Service
	YaruzRepository yarus_platform.IPlatform
}

type Auth struct {
	//SessionRepository auth.SessionRepository
	//TokenRepository   auth.TokenRepository
	//Service           auth.Service
}

// Domain is a Domain Layer Entry Point
type Domain struct {
	User                          user.IService
	userRepository                user.Repository
	Auth                          auth.Service
	SessionRepository             session.Repository
	SnAccount                     sn_account.IService
	snAccountRepository           sn_account.Repository
	Advertiser                    advertiser.IService
	advertiserRepository          advertiser.Repository
	AdvertisingCampaign           advertising_campaign.IService
	advertisingCampaignRepository advertising_campaign.Repository
	Offer                         offer.IService
	offerRepository               offer.Repository
}

// New func is a constructor for the App
func New(ctx context.Context, cfg config.Configuration) *App {
	logger, err := log.New(cfg.Log)
	if err != nil {
		golog.Fatal(err)
	}

	infra, err := NewInfra(ctx, logger, cfg)
	if err != nil {
		golog.Fatal(err)
	}

	app := &App{
		Cfg:   cfg,
		Infra: infra,
	}

	err = app.Init(ctx)
	if err != nil {
		golog.Fatal(err)
	}

	return app
}
func NewInfra(ctx context.Context, logger log.ILogger, cfg config.Configuration) (*Infrastructure, error) {
	IdentityDB, err := minipkg_gorm.New(logger, cfg.DB.Identity)
	if err != nil {
		return nil, err
	}

	rDB, err := redis.New(cfg.DB.Redis)
	if err != nil {
		return nil, err
	}

	yaruzRepository, err := yarus_platform.NewPlatform(ctx, cfg.YaruzConfig())
	if err != nil {
		return nil, err
	}

	return &Infrastructure{
		Logger:          logger,
		IdentityDB:      IdentityDB,
		Redis:           rDB,
		YaruzRepository: yaruzRepository,
	}, nil
}

func (app *App) Init(ctx context.Context) (err error) {
	if err := app.SetupRepositories(); err != nil {
		return err
	}

	return app.SetupServices(ctx)
}

func (app *App) SetupRepositories() (err error) {
	var ok bool

	app.Domain.SessionRepository, err = redisrepo.NewSessionRepository(app.Infra.Redis, app.Cfg.Auth.SessionlifeTime)
	if err != nil {
		golog.Fatalf("Can not get session repository, error happened: %v", err)
	}

	userRepo, err := yaruzplatform.GetRepository(app.Infra.Logger, app.Infra.YaruzRepository, user.EntityType)
	if err != nil {
		golog.Fatalf("Can not get yaruz repository for entity %q, error happened: %v", user.EntityType, err)
	}

	app.Domain.userRepository, ok = userRepo.(user.Repository)
	if !ok {
		return errors.Errorf("Can not cast yaruz repository for entity %q to %vRepository. Repo: %v", user.EntityType, user.EntityType, userRepo)
	}

	snAccountRepo, err := yaruzplatform.GetRepository(app.Infra.Logger, app.Infra.YaruzRepository, sn_account.EntityType)
	if err != nil {
		golog.Fatalf("Can not get yaruz repository for entity %q, error happened: %v", sn_account.EntityType, err)
	}

	app.Domain.snAccountRepository, ok = snAccountRepo.(sn_account.Repository)
	if !ok {
		return errors.Errorf("Can not cast yaruz repository for entity %q to %vRepository. Repo: %v", sn_account.EntityType, sn_account.EntityType, snAccountRepo)
	}

	advertiserRepo, err := yaruzplatform.GetRepository(app.Infra.Logger, app.Infra.YaruzRepository, advertiser.EntityType)
	if err != nil {
		golog.Fatalf("Can not get yaruz repository for entity %q, error happened: %v", advertiser.EntityType, err)
	}

	app.Domain.advertiserRepository, ok = advertiserRepo.(advertiser.Repository)
	if !ok {
		return errors.Errorf("Can not cast yaruz repository for entity %q to %vRepository. Repo: %v", advertiser.EntityType, advertiser.EntityType, advertiserRepo)
	}

	advertisingCampaignRepo, err := yaruzplatform.GetRepository(app.Infra.Logger, app.Infra.YaruzRepository, advertising_campaign.EntityType)
	if err != nil {
		golog.Fatalf("Can not get yaruz repository for entity %q, error happened: %v", advertising_campaign.EntityType, err)
	}

	app.Domain.advertisingCampaignRepository, ok = advertisingCampaignRepo.(advertising_campaign.Repository)
	if !ok {
		return errors.Errorf("Can not cast yaruz repository for entity %q to %vRepository. Repo: %v", advertising_campaign.EntityType, advertising_campaign.EntityType, advertisingCampaignRepo)
	}

	offerRepo, err := yaruzplatform.GetRepository(app.Infra.Logger, app.Infra.YaruzRepository, offer.EntityType)
	if err != nil {
		golog.Fatalf("Can not get yaruz repository for entity %q, error happened: %v", offer.EntityType, err)
	}

	app.Domain.offerRepository, ok = offerRepo.(offer.Repository)
	if !ok {
		return errors.Errorf("Can not cast yaruz repository for entity %q to %vRepository. Repo: %v", offer.EntityType, offer.EntityType, offerRepo)
	}

	//if app.Auth.SessionRepository, err = redisrep.NewSessionRepository(app.Infra.Redis, app.Cfg.SessionLifeTime, app.Domain.User.Repository); err != nil {
	//	return errors.Errorf("Can not get new SessionRepository err: %v", err)
	//}
	//app.Auth.TokenRepository = jwt.NewRepository()

	app.Infra.Cache = cache.NewService(app.Infra.Redis, app.Cfg.CacheLifeTime)

	return nil
}

func (app *App) SetupServices(ctx context.Context) error {
	var err error
	app.Domain.User = user.NewService(app.Infra.Logger, app.Domain.userRepository)
	app.Domain.Auth, err = auth.NewService(ctx, app.Infra.Logger, app.Cfg.Auth, app.Domain.User, app.Domain.SessionRepository, app.Infra.YaruzRepository.ReferenceSubsystem().TextLang)
	if err != nil {
		return err
	}
	app.Domain.SnAccount = sn_account.NewService(app.Infra.Logger, app.Domain.snAccountRepository)
	app.Domain.Advertiser = advertiser.NewService(app.Infra.Logger, app.Domain.advertiserRepository)
	app.Domain.AdvertisingCampaign = advertising_campaign.NewService(app.Infra.Logger, app.Domain.advertisingCampaignRepository)
	app.Domain.Offer = offer.NewService(app.Infra.Logger, app.Domain.offerRepository)

	return nil
}

// Run is func to run the App
func (app *App) Run() error {
	return nil
}

func (app *App) Stop() error {
	errRedis := app.Infra.Redis.Close()
	errDB01 := app.Infra.IdentityDB.Close()
	errDB02 := app.Infra.YaruzRepository.Stop()

	switch {
	case errDB01 != nil:
		return errors.Wrapf(apperror.ErrInternal, "db close error: %v", errDB01)
	case errDB02 != nil:
		return errors.Wrapf(apperror.ErrInternal, "yarus repository close error: %v", errDB02)
	case errRedis != nil:
		return errors.Wrapf(apperror.ErrInternal, "redis close error: %v", errRedis)
	}

	return nil
}
