package data

import (
	golog "log"

	"github.com/yaruz/app/pkg/yarus_platform/search"

	"github.com/yaruz/app/pkg/yarus_platform/reference"

	"github.com/yaruz/app/pkg/yarus_platform/infrastructure"

	"github.com/yaruz/app/pkg/yarus_platform/data/infrastructure/repository/gorm"

	"github.com/minipkg/log"
	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"
)

type DataSubsystem struct {
	reference         *reference.ReferenceSubsystem
	search            *search.SearchSubsystem
	Entity            DataDomainEntity
	ValueRepositories entity.ValueRepositories
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
	return d, nil
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

	repo, err = gorm.GetRepository(infra.Logger, infra.DataDB, text_value.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", text_value.EntityName, err)
	}
	d.ValueRepositories.Text, ok = repo.(text_value.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", text_value.EntityName, text_value.EntityName, repo)
	}

	return nil
}

func (d *DataSubsystem) setupServices(logger log.ILogger) {
	d.Entity.Service = entity.NewService(logger, d.Entity.Repository, d.reference, d.search)
}
