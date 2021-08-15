package yarus_platform

import (
	"github.com/minipkg/log"
	"github.com/yaruz/app/pkg/yarus_platform/data"
	"github.com/yaruz/app/pkg/yarus_platform/infrastructure"
	"github.com/yaruz/app/pkg/yarus_platform/reference"
	"github.com/yaruz/app/pkg/yarus_platform/search"

	"github.com/yaruz/app/pkg/yarus_platform/config"
)

type IPlatform interface {
	Stop() error
	DataSubsystem() *data.DataSubsystem
	ReferenceSubsystem() *reference.ReferenceSubsystem
	SearchSubsystem() *search.SearchSubsystem
}

type Platform struct {
	cfg                config.Configuration
	dataSubsystem      *data.DataSubsystem
	referenceSubsystem *reference.ReferenceSubsystem
	searchSubsystem    *search.SearchSubsystem
	infra              *infrastructure.Infrastructure
}

func NewPlatform(cfg config.Configuration) (*Platform, error) {
	logger, err := log.New(cfg.Infra.Log)
	if err != nil {
		return nil, err
	}

	infra, err := infrastructure.NewInfra(logger, cfg.Infra)
	if err != nil {
		return nil, err
	}

	data, err := data.NewDataSubsystem(infra)
	if err != nil {
		return nil, err
	}

	reference, err := reference.NewReferenceSubsystem(infra)
	if err != nil {
		return nil, err
	}

	search, err := search.NewSearchSubsystem(infra)
	if err != nil {
		return nil, err
	}

	return &Platform{
		cfg:                cfg,
		dataSubsystem:      data,
		referenceSubsystem: reference,
		searchSubsystem:    search,
		infra:              infra,
	}, nil
}

func (r *Platform) DataSubsystem() *data.DataSubsystem {
	return r.dataSubsystem
}

func (r *Platform) ReferenceSubsystem() *reference.ReferenceSubsystem {
	return r.referenceSubsystem
}

func (r *Platform) SearchSubsystem() *search.SearchSubsystem {
	return r.searchSubsystem
}

func (r *Platform) Stop() error {
	return r.infra.Stop()
}
