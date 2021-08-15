package search

import (
	"github.com/minipkg/log"
	"github.com/yaruz/app/pkg/yarus_platform/infrastructure"
)

type SearchSubsystem struct {
}

func NewSearchSubsystem(infra *infrastructure.Infrastructure) (*SearchSubsystem, error) {
	d := &SearchSubsystem{}
	if err := d.setupRepositories(infra); err != nil {
		return nil, err
	}
	d.setupServices(infra.Logger)
	return d, nil
}

func (d *SearchSubsystem) setupRepositories(infra *infrastructure.Infrastructure) (err error) {
	return nil
}

func (d *SearchSubsystem) setupServices(logger log.ILogger) {

}
