package yarus_platform

import (
	"github.com/minipkg/log"
)

type SearchDomain struct {
}

func newSearchDomain(logger log.ILogger, infra *infrastructure) (*SearchDomain, error) {
	d := &SearchDomain{}
	if err := d.setupRepositories(logger, infra); err != nil {
		return nil, err
	}
	d.setupServices(logger)
	return d, nil
}

func (d *SearchDomain) setupRepositories(logger log.ILogger, infra *infrastructure) (err error) {
	return nil
}

func (d *SearchDomain) setupServices(logger log.ILogger) {

}
