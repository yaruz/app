package yarus_platform

import (
	golog "log"

	"github.com/yaruz/app/pkg/yarus_platform/data/infrastructure/repository/gorm"

	"github.com/minipkg/log"
	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"
)

type DataDomain struct {
	Entity    DataDomainEntity
	TextValue DataDomainTextValue
}

type DataDomainEntity struct {
	Service    entity.IService
	Repository entity.Repository
}

type DataDomainTextValue struct {
	Service    text_value.IService
	Repository text_value.Repository
}

func newDataDomain(infra *infrastructure) (*DataDomain, error) {
	d := &DataDomain{}
	if err := d.setupRepositories(infra); err != nil {
		return nil, err
	}
	d.setupServices(infra.Logger)
	return d, nil
}

func (d *DataDomain) setupRepositories(infra *infrastructure) (err error) {
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
	d.TextValue.Repository, ok = repo.(text_value.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", text_value.EntityName, text_value.EntityName, repo)
	}

	return nil
}

func (d *DataDomain) setupServices(logger log.ILogger) {
	d.Entity.Service = entity.NewService(logger, d.Entity.Repository)
	d.TextValue.Service = text_value.NewService(logger, d.TextValue.Repository)
}
