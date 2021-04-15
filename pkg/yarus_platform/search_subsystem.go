package yarus_platform

import (
	"github.com/minipkg/log"
)

type SearchDomain struct {
}

func newSearchDomain(infra *infrastructure) (*SearchDomain, error) {
	d := &SearchDomain{}
	if err := d.setupRepositories(infra); err != nil {
		return nil, err
	}
	d.setupServices(infra.Logger)
	return d, nil
}

func (d *SearchDomain) setupRepositories(infra *infrastructure) (err error) {
	return nil
}

func (d *SearchDomain) setupServices(logger log.ILogger) {

}
