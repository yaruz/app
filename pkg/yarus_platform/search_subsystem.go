package yarus_platform

import (
	"github.com/minipkg/log"
)

type SearchSubsystem struct {
}

func newSearchSubsystem(infra *infrastructure) (*SearchSubsystem, error) {
	d := &SearchSubsystem{}
	if err := d.setupRepositories(infra); err != nil {
		return nil, err
	}
	d.setupServices(infra.Logger)
	return d, nil
}

func (d *SearchSubsystem) setupRepositories(infra *infrastructure) (err error) {
	return nil
}

func (d *SearchSubsystem) setupServices(logger log.ILogger) {

}
