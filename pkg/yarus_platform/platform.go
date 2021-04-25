package yarus_platform

import (
	"github.com/minipkg/log"

	"github.com/yaruz/app/pkg/yarus_platform/config"
)

type IPlatform interface {
	Stop() error
	Data() *DataDomain
	Reference() *ReferenceDomain
	Search() *SearchDomain
}

type Platform struct {
	cfg             config.Configuration
	DataDomain      *DataDomain
	ReferenceDomain *ReferenceDomain
	SearchDomain    *SearchDomain
	infra           *infrastructure
}

func NewPlatform(cfg config.Configuration) (*Platform, error) {
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

	reference, err := newReferenceDomain(infra)
	if err != nil {
		return nil, err
	}

	search, err := newSearchDomain(infra)
	if err != nil {
		return nil, err
	}

	return &Platform{
		cfg:             cfg,
		DataDomain:      data,
		ReferenceDomain: reference,
		SearchDomain:    search,
		infra:           infra,
	}, nil
}

func (r *Platform) Data() *DataDomain {
	return r.DataDomain
}

func (r *Platform) Reference() *ReferenceDomain {
	return r.ReferenceDomain
}

func (r *Platform) Search() *SearchDomain {
	return r.SearchDomain
}

func (r *Platform) Stop() error {
	return r.infra.Stop()
}
