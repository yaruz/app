package infrastructure

import (
	"context"
	golog "log"

	"github.com/minipkg/db/redis"
	"github.com/minipkg/log"

	"github.com/yaruz/app/pkg/yarus_platform"
	yaruz_config "github.com/yaruz/app/pkg/yarus_platform/config"

	redis_repo "github.com/yaruz/app/internal/infrastructure/repository/redis"
)

type Infrastructure struct {
	Logger          log.Logger
	Redis           redis_repo.Repository
	YaruzRepository yarus_platform.IPlatform
}

func New(ctx context.Context, cfg *Configuration, yaruzMetadata *yaruz_config.Metadata) (*Infrastructure, error) {
	logger, err := log.New(cfg.Log)
	if err != nil {
		golog.Fatal(err)
	}

	rDB, err := redis.New(cfg.DB.Redis)
	if err != nil {
		return nil, err
	}

	yaruzRepository, err := yarus_platform.NewPlatform(ctx, cfg.YaruzConfig(yaruzMetadata))
	if err != nil {
		return nil, err
	}

	return &Infrastructure{
		Logger:          logger,
		Redis:           redis_repo.NewRepository(rDB),
		YaruzRepository: yaruzRepository,
	}, nil
}
