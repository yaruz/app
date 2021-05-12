package yarus_platform

import (
	"github.com/minipkg/log"

	"github.com/yaruz/app/pkg/yarus_platform/config"
)

type IPlatform interface {
	Stop() error
	DataSubsystem() *DataSubsystem
	ReferenceSubsystem() *ReferenceSubsystem
	SearchSubsystem() *SearchSubsystem
}

type Platform struct {
	cfg                config.Configuration
	dataSubsystem      *DataSubsystem
	referenceSubsystem *ReferenceSubsystem
	searchSubsystem    *SearchSubsystem
	infra              *infrastructure
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

	data, err := newDataSubsystem(infra)
	if err != nil {
		return nil, err
	}

	reference, err := newReferenceSubsystem(infra)
	if err != nil {
		return nil, err
	}

	search, err := newSearchSubsystem(infra)
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

func (r *Platform) DataSubsystem() *DataSubsystem {
	return r.dataSubsystem
}

func (r *Platform) ReferenceSubsystem() *ReferenceSubsystem {
	return r.referenceSubsystem
}

func (r *Platform) SearchSubsystem() *SearchSubsystem {
	return r.searchSubsystem
}

func (r *Platform) Stop() error {
	return r.infra.Stop()
}
