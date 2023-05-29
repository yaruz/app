package app

import (
	"context"
	golog "log"

	"github.com/pkg/errors"
	redisrepo "github.com/yaruz/app/internal/infrastructure/repository/redis"

	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/internal/pkg/auth"
	"github.com/yaruz/app/internal/pkg/config"
	"github.com/yaruz/app/internal/pkg/jwt"
	"github.com/yaruz/app/internal/pkg/session"
	"github.com/yaruz/app/internal/pkg/socnets/tg"

	"github.com/yaruz/app/internal/domain/advertiser"
	"github.com/yaruz/app/internal/domain/advertising_campaign"
	"github.com/yaruz/app/internal/domain/offer"
	"github.com/yaruz/app/internal/domain/tg_account"
	"github.com/yaruz/app/internal/domain/user"
	"github.com/yaruz/app/internal/infrastructure"
	"github.com/yaruz/app/internal/infrastructure/repository/yaruzplatform"
)

// App struct is the common part of all applications
type App struct {
	Cfg    config.Configuration
	Domain Domain
	Auth   Auth
	Infra  *infrastructure.Infrastructure
}

type Auth struct {
	//sessionRepository auth.sessionRepository
	//TokenRepository   auth.TokenRepository
	//Service           auth.Service
}

// Domain is a Domain Layer Entry Point
type Domain struct {
	User                          user.IService
	userRepository                user.Repository
	Auth                          auth.Service
	sessionRepository             session.Repository
	jwtRepository                 auth.TokenRepository
	Tg                            tg.Service
	TgAccount                     tg_account.IService
	tgAccountRepository           tg_account.Repository
	Advertiser                    advertiser.IService
	advertiserRepository          advertiser.Repository
	AdvertisingCampaign           advertising_campaign.IService
	advertisingCampaignRepository advertising_campaign.Repository
	Offer                         offer.IService
	offerRepository               offer.Repository
}

// New func is a constructor for the App
func New(ctx context.Context, cfg config.Configuration) *App {
	infra, err := infrastructure.New(ctx, &cfg.Infrastructure, &cfg.YaruzMetadata)
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

func (app *App) Init(ctx context.Context) (err error) {
	if err := app.SetupRepositories(); err != nil {
		return err
	}

	return app.SetupServices(ctx)
}

func (app *App) SetupRepositories() (err error) {
	var ok bool
	app.Domain.jwtRepository = jwt.NewRepository(app.Cfg.Auth.JWTSigningKey, app.Cfg.Auth.JWTExpirationInHours)

	app.Domain.sessionRepository, err = redisrepo.NewSessionRepository(app.Infra.Redis, app.Cfg.Auth.SessionlifeTimeInHours)
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

	tgAccountRepo, err := yaruzplatform.GetRepository(app.Infra.Logger, app.Infra.YaruzRepository, tg_account.EntityType)
	if err != nil {
		golog.Fatalf("Can not get yaruz repository for entity %q, error happened: %v", tg_account.EntityType, err)
	}

	app.Domain.tgAccountRepository, ok = tgAccountRepo.(tg_account.Repository)
	if !ok {
		return errors.Errorf("Can not cast yaruz repository for entity %q to %vRepository. Repo: %v", tg_account.EntityType, tg_account.EntityType, tgAccountRepo)
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

	//if app.Auth.sessionRepository, err = redisrep.NewSessionRepository(app.Infra.Redis, app.Cfg.SessionLifeTime, app.Domain.User.Repository); err != nil {
	//	return errors.Errorf("Can not get new sessionRepository err: %v", err)
	//}
	//app.Auth.TokenRepository = jwt.NewRepository()

	//app.Infra.Cache = cache.NewService(app.Infra.Redis, app.Cfg.CacheLifeTime)

	return nil
}

func (app *App) SetupServices(ctx context.Context) error {
	var err error
	app.Domain.User = user.NewService(app.Infra.Logger, app.Domain.userRepository)
	app.Domain.Auth, err = auth.NewService(ctx, app.Infra.Logger, app.Cfg.Auth, app.Domain.jwtRepository, app.Domain.sessionRepository, app.Domain.User, app.Infra.YaruzRepository.ReferenceSubsystem().TextLang)
	if err != nil {
		return err
	}
	app.Domain.TgAccount = tg_account.NewService(app.Infra.Logger, app.Domain.tgAccountRepository)
	app.Domain.Advertiser = advertiser.NewService(app.Infra.Logger, app.Domain.advertiserRepository)
	app.Domain.AdvertisingCampaign = advertising_campaign.NewService(app.Infra.Logger, app.Domain.advertisingCampaignRepository)
	app.Domain.Offer = offer.NewService(app.Infra.Logger, app.Domain.offerRepository)

	app.Domain.Tg = tg.New(app.Cfg.Socnets.Telegram, app.Infra.Logger, redisrepo.NewTgSessionRepository(app.Infra.Redis))

	return nil
}

// Run is func to run the App
func (app *App) Run() error {
	return nil
}

func (app *App) Stop() error {
	errRedis := app.Infra.Redis.Close()
	err := app.Infra.YaruzRepository.Stop()

	switch {
	case err != nil:
		return errors.Wrapf(apperror.ErrInternal, "yarus repository close error: %v", err)
	case errRedis != nil:
		return errors.Wrapf(apperror.ErrInternal, "redis close error: %v", errRedis)
	}

	return nil
}
