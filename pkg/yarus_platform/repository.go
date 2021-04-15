package yarus_platform

import (
	"github.com/minipkg/log"

	"github.com/yaruz/app/pkg/yarus_platform/config"
)

type IRepository interface {
	Stop() error
}

type Repository struct {
	cfg    config.Configuration
	Data   *DataDomain
	Search *SearchDomain
	infra  *infrastructure
}

func NewRepository(cfg config.Configuration) (*Repository, error) {
	logger, err := log.New(cfg.Infra.Log)
	if err != nil {
		return nil, err
	}

	infra, err := newInfra(logger, cfg.Infra)
	if err != nil {
		return nil, err
	}

	data, err := newDataDomain(infra)
	if err != nil {
		return nil, err
	}

	search, err := newSearchDomain(infra)
	if err != nil {
		return nil, err
	}

	return &Repository{
		cfg:    cfg,
		Data:   data,
		Search: search,
		infra:  infra,
	}, nil
}

func (r *Repository) Stop() error {
	return r.infra.Stop()
}
