package yarus_platform

import (
	golog "log"

	"github.com/minipkg/log"

	"github.com/yaruz/app/pkg/yarus_platform/config"
)

type Repository struct {
	cfg    config.Configuration
	Logger log.ILogger
	Data   DataDomain
	Search SearchDomain
	infra  *infrastructure
}

func NewRepository(cfg config.Configuration) *Repository {
	logger, err := log.New(cfg.Log)
	if err != nil {
		golog.Fatal(err)
	}

	infra, err := newInfra(cfg.Infra, logger)
	if err != nil {
		golog.Fatal(err)
	}

	return &Repository{
		cfg:    cfg,
		Logger: logger,
		Data:   DataDomain{},
		Search: SearchDomain{},
		infra:  infra,
	}
}
