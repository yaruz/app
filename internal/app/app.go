package app

import (
	golog "log"

	"github.com/yaruz/app/internal/infrastructure/repository/yaruzplatform"

	"github.com/yaruz/app/pkg/yarus_platform"

	"github.com/pkg/errors"

	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/internal/pkg/auth"
	"github.com/yaruz/app/internal/pkg/config"
	"github.com/yaruz/app/internal/pkg/jwt"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/db/redis"
	"github.com/minipkg/db/redis/cache"
	"github.com/minipkg/log"

	"github.com/yaruz/app/internal/domain/task"
	"github.com/yaruz/app/internal/domain/user"
	gormrep "github.com/yaruz/app/internal/infrastructure/repository/gorm"
	redisrep "github.com/yaruz/app/internal/infrastructure/repository/redis"
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
	yaruzRepository yarus_platform.IRepository
}

type Auth struct {
	SessionRepository auth.SessionRepository
	TokenRepository   auth.TokenRepository
	Service           auth.Service
}

// Domain is a Domain Layer Entry Point
type Domain struct {
	User DomainUser
	//	Example
	Task DomainTask
}

type DomainUser struct {
	Repository user.Repository
	Service    user.IService
}

type DomainTask struct {
	Repository task.Repository
	Service    task.IService
}

// New func is a constructor for the App
func New(cfg config.Configuration) *App {
	logger, err := log.New(cfg.Log)
	if err != nil {
		golog.Fatal(err)
	}

	infra, err := NewInfra(logger, cfg)
	if err != nil {
		golog.Fatal(err)
	}

	app := &App{
		Cfg:   cfg,
		Infra: infra,
	}

	err = app.Init()
	if err != nil {
		golog.Fatal(err)
	}

	return app
}
func NewInfra(logger log.ILogger, cfg config.Configuration) (*Infrastructure, error) {
	IdentityDB, err := minipkg_gorm.New(logger, cfg.DB.Identity)
	if err != nil {
		return nil, err
	}

	rDB, err := redis.New(cfg.DB.Redis)
	if err != nil {
		return nil, err
	}

	yaruzRepository, err := yarus_platform.NewRepository(cfg.YaruzPlatformConfig())
	if err != nil {
		return nil, err
	}

	return &Infrastructure{
		Logger:          logger,
		IdentityDB:      IdentityDB,
		Redis:           rDB,
		yaruzRepository: yaruzRepository,
	}, nil
}

func (app *App) Init() (err error) {
	if err := app.SetupRepositories(); err != nil {
		return err
	}
	app.SetupServices()
	return nil
}

func (app *App) SetupRepositories() (err error) {
	var ok bool

	gormRepo, err := gormrep.GetRepository(app.Infra.Logger, app.Infra.IdentityDB, user.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", user.EntityName, err)
	}

	app.Domain.User.Repository, ok = gormRepo.(user.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", user.EntityName, user.EntityName, gormRepo)
	}
	//	Example
	yarRepo, err := yaruzplatform.GetRepository(app.Infra.Logger, app.Infra.yaruzRepository, task.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", task.EntityName, err)
	}

	app.Domain.Task.Repository, ok = yarRepo.(task.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", task.EntityName, task.EntityName, yarRepo)
	}

	if app.Auth.SessionRepository, err = redisrep.NewSessionRepository(app.Infra.Redis, app.Cfg.SessionLifeTime, app.Domain.User.Repository); err != nil {
		return errors.Errorf("Can not get new SessionRepository err: %v", err)
	}
	app.Auth.TokenRepository = jwt.NewRepository()

	app.Infra.Cache = cache.NewService(app.Infra.Redis, app.Cfg.CacheLifeTime)

	return nil
}

func (app *App) SetupServices() {
	app.Domain.User.Service = user.NewService(app.Infra.Logger, app.Domain.User.Repository)
	app.Auth.Service = auth.NewService(app.Cfg.JWTSigningKey, app.Cfg.JWTExpiration, app.Domain.User.Service, app.Infra.Logger, app.Auth.SessionRepository, app.Auth.TokenRepository)
	//	Example
	app.Domain.Task.Service = task.NewService(app.Infra.Logger, app.Domain.Task.Repository)
}

// Run is func to run the App
func (app *App) Run() error {
	return nil
}

func (app *App) Stop() error {
	errRedis := app.Infra.Redis.Close()
	errDB01 := app.Infra.IdentityDB.DB().Close()
	errDB02 := app.Infra.yaruzRepository.Stop()

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
