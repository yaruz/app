package yarus_platform

import (
	"context"

	"github.com/minipkg/log"
	"github.com/yaruz/app/pkg/yarus_platform/config"
	"github.com/yaruz/app/pkg/yarus_platform/data"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/infrastructure"
	"github.com/yaruz/app/pkg/yarus_platform/reference"
)

type IPlatform interface {
	Stop() error
	DataSubsystem() *data.DataSubsystem
	ReferenceSubsystem() *reference.ReferenceSubsystem
	//SearchSubsystem() *search.SearchSubsystem
	AutoMigrate(metadata *config.Metadata) error
}

type Platform struct {
	cfg                config.Configuration
	dataSubsystem      *data.DataSubsystem
	referenceSubsystem *reference.ReferenceSubsystem
	//searchSubsystem    *search.SearchSubsystem
	infra *infrastructure.Infrastructure
}

func NewPlatform(ctx context.Context, cfg config.Configuration) (*Platform, error) {
	logger, err := log.New(cfg.Infrastructure.Log)
	if err != nil {
		return nil, err
	}

	infra, err := infrastructure.NewInfrastructure(ctx, logger, cfg.Infrastructure, entity.New())
	if err != nil {
		return nil, err
	}

	reference, err := reference.NewReferenceSubsystem(infra, cfg.Metadata)
	if err != nil {
		return nil, err
	}

	//search, err := search.NewSearchSubsystem(infra)
	//if err != nil {
	//	return nil, err
	//}

	data, err := data.NewDataSubsystem(infra, reference)
	if err != nil {
		return nil, err
	}

	return &Platform{
		cfg:                cfg,
		dataSubsystem:      data,
		referenceSubsystem: reference,
		//searchSubsystem:    search,
		infra: infra,
	}, nil
}

func (r *Platform) DataSubsystem() *data.DataSubsystem {
	return r.dataSubsystem
}

func (r *Platform) ReferenceSubsystem() *reference.ReferenceSubsystem {
	return r.referenceSubsystem
}

//
//func (r *Platform) SearchSubsystem() *search.SearchSubsystem {
//	return r.searchSubsystem
//}

func (r *Platform) Stop() error {
	return r.infra.Stop()
}

func (r *Platform) AutoMigrate(metadata *config.Metadata) error {
	return nil
}
