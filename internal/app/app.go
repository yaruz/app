package app

import (
	golog "log"

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
	pgrep "github.com/yaruz/app/internal/infrastructure/repository/gorm"
	redisrep "github.com/yaruz/app/internal/infrastructure/repository/redis"
)

// App struct is the common part of all applications
type App struct {
	Cfg        config.Configuration
	Logger     log.ILogger
	IdentityDB minipkg_gorm.IDB
	DataDB     minipkg_gorm.IDB
	SearchDB   minipkg_gorm.IDB
	Redis      redis.IDB
	Domain     Domain
	Auth       Auth
	Cache      cache.Service
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

	IdentityDB, err := minipkg_gorm.New(cfg.DB.Identity, logger)
	if err != nil {
		golog.Fatal(err)
	}

	DataDB, err := minipkg_gorm.New(cfg.DB.Data, logger)
	if err != nil {
		golog.Fatal(err)
	}

	SearchDB, err := minipkg_gorm.New(cfg.DB.Search, logger)
	if err != nil {
		golog.Fatal(err)
	}

	rDB, err := redis.New(cfg.DB.Redis)
	if err != nil {
		golog.Fatal(err)
	}

	app := &App{
		Cfg:        cfg,
		Logger:     logger,
		IdentityDB: IdentityDB,
		DataDB:     DataDB,
		SearchDB:   SearchDB,
		Redis:      rDB,
	}

	err = app.Init()
	if err != nil {
		golog.Fatal(err)
	}

	return app
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

	app.Domain.User.Repository, ok = app.getPgRepo(app.IdentityDB, user.EntityName).(user.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", user.EntityName, user.EntityName, app.getPgRepo(app.IdentityDB, user.EntityName))
	}
	//	CarCatalog
	app.Domain.Mark.Repository, ok = app.getPgRepo(app.DataDB, mark.EntityName).(mark.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", mark.EntityName, mark.EntityName, app.getPgRepo(app.DataDB, mark.EntityName))
	}

	app.Domain.Task.Repository, ok = app.getPgRepo(app.DataDB, task.EntityName).(task.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", task.EntityName, task.EntityName, app.getPgRepo(app.DataDB, task.EntityName))
	}

	if app.Auth.SessionRepository, err = redisrep.NewSessionRepository(app.Redis, app.Cfg.SessionLifeTime, app.Domain.User.Repository); err != nil {
		return errors.Errorf("Can not get new SessionRepository err: %v", err)
	}
	app.Auth.TokenRepository = jwt.NewRepository()

	app.Cache = cache.NewService(app.Redis, app.Cfg.CacheLifeTime)

	return nil
}

func (app *App) SetupServices() {
	app.Domain.User.Service = user.NewService(app.Logger, app.Domain.User.Repository)
	app.Auth.Service = auth.NewService(app.Cfg.JWTSigningKey, app.Cfg.JWTExpiration, app.Domain.User.Service, app.Logger, app.Auth.SessionRepository, app.Auth.TokenRepository)
	//	CarCatalog
	app.Domain.Mark.Service = mark.NewService(app.Logger, app.Domain.Mark.Repository)
	app.Domain.Task.Service = task.NewService(app.Logger, app.Domain.Task.Repository)
}

// Run is func to run the App
func (app *App) Run() error {
	return nil
}

func (app *App) getPgRepo(dbase minipkg_gorm.IDB, entityName string) (repo pgrep.IRepository) {
	var err error

	if repo, err = pgrep.GetRepository(app.Logger, dbase, entityName); err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entityName, err)
	}
	return repo
}

func (app *App) Stop() error {
	errRedis := app.Redis.Close()
	errDB01 := app.IdentityDB.DB().Close()
	errDB02 := app.DataDB.DB().Close()
	errDB03 := app.SearchDB.DB().Close()

	switch {
	case errDB01 != nil || errDB02 != nil || errDB03 != nil:
		return errors.Wrapf(apperror.ErrInternal, "db close error: %v", errDB01)
	case errRedis != nil:
		return errors.Wrapf(apperror.ErrInternal, "redis close error: %v", errRedis)
	}

	return nil
}
