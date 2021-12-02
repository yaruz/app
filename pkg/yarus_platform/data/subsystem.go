package data

import (
	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/log"
	"github.com/pkg/errors"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/bool_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/date_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/float_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/int_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/time_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/infrastructure/repository/gorm"
	"github.com/yaruz/app/pkg/yarus_platform/infrastructure"
	"github.com/yaruz/app/pkg/yarus_platform/reference"
)

type DataSubsystem struct {
	reference        *reference.ReferenceSubsystem
	Entity           entity.IService
	entityRepository entity.Repository
	search           entity.SearchService
	mapReducer       gorm.IMapReducer
}

func NewDataSubsystem(infra *infrastructure.Infrastructure, reference *reference.ReferenceSubsystem) (*DataSubsystem, error) {
	d := &DataSubsystem{
		reference: reference,
	}
	if err := d.setupRepositories(infra); err != nil {
		return nil, err
	}
	d.setupServices(infra.Logger)

	if err := d.autoMigrate(infra.DataSharding); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *DataSubsystem) autoMigrate(sharding infrastructure.Sharding) error {
	if sharding.IsAutoMigrate {
		return sharding.ApplyFunc2DBs(func(db minipkg_gorm.IDB) error {
			return db.DB().AutoMigrate(
				&entity.Entity{},
				&bool_value.BoolValue{},
				&int_value.IntValue{},
				&float_value.FloatValue{},
				&date_value.DateValue{},
				&time_value.TimeValue{},
				&text_value.TextValue{},
			)
		})
	}
	return nil
}

func (d *DataSubsystem) setupRepositories(infra *infrastructure.Infrastructure) (err error) {
	var ok bool

	mr := gorm.NewMapReducer(infra.Logger, d.reference.EntityType, infra.DataSharding)

	repo, err := gorm.GetRepository(infra.Logger, mr, entity.EntityName, d.reference.TextLang)
	if err != nil {
		return errors.Errorf("Can not get db repository for entity %q, error happened: %v", entity.EntityName, err)
	}
	d.entityRepository, ok = repo.(entity.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity.EntityName, entity.EntityName, repo)
	}

	d.search = gorm.NewSearchService(infra.Logger, mr, d.reference.EntityType, d.reference.Property, d.reference.TextLang)

	return nil
}

func (d *DataSubsystem) setupServices(logger log.ILogger) {
	d.Entity = entity.NewService(logger, d.entityRepository, d.reference, d.search)
}
