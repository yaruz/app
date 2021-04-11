package yarus_platform

import (
	golog "log"

	"github.com/minipkg/log"

	"github.com/yaruz/app/pkg/yarus_platform/config"
)

type Repository struct {
	cfg    config.Configuration
	Logger log.ILogger
	Data   *DataDomain
	Search *SearchDomain
	infra  *infrastructure
}

func NewRepository(cfg config.Configuration) *Repository {
	logger, err := log.New(cfg.Log)
	if err != nil {
		golog.Fatal(err)
	}

	infra, err := newInfra(logger, cfg.Infra)
	if err != nil {
		golog.Fatal(err)
	}

	data, err := newDataDomain(logger, infra)
	if err != nil {
		golog.Fatal(err)
	}

	search, err := newSearchDomain(logger, infra)
	if err != nil {
		golog.Fatal(err)
	}

	return &Repository{
		cfg:    cfg,
		Logger: logger,
		Data:   data,
		Search: search,
		infra:  infra,
	}
}
