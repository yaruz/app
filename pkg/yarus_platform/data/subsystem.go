package data

import (
	golog "log"

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
	"github.com/yaruz/app/pkg/yarus_platform/search"
)

type DataSubsystem struct {
	reference *reference.ReferenceSubsystem
	search    *search.SearchSubsystem
	Entity    DataDomainEntity
}

type DataDomainEntity struct {
	Service    entity.IService
	Repository entity.Repository
}

func NewDataSubsystem(infra *infrastructure.Infrastructure, reference *reference.ReferenceSubsystem, search *search.SearchSubsystem) (*DataSubsystem, error) {
	d := &DataSubsystem{
		reference: reference,
		search:    search,
	}
	if err := d.setupRepositories(infra); err != nil {
		return nil, err
	}
	d.setupServices(infra.Logger)

	if err := d.autoMigrate(infra.DataDB); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *DataSubsystem) autoMigrate(db minipkg_gorm.IDB) error {
	if db.IsAutoMigrate() {

		err := db.DB().AutoMigrate(
			&bool_value.BoolValue{},
			&int_value.IntValue{},
			&float_value.FloatValue{},
			&date_value.DateValue{},
			&time_value.TimeValue{},
			&text_value.TextValue{},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DataSubsystem) setupRepositories(infra *infrastructure.Infrastructure) (err error) {
	var ok bool

	repo, err := gorm.GetRepository(infra.Logger, infra.DataDB, entity.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entity.EntityName, err)
	}
	d.Entity.Repository, ok = repo.(entity.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity.EntityName, entity.EntityName, repo)
	}

	return nil
}

func (d *DataSubsystem) setupServices(logger log.ILogger) {
	d.Entity.Service = entity.NewService(logger, d.Entity.Repository, d.reference, d.search)
}
