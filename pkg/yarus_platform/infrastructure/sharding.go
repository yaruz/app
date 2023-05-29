package infrastructure

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"

	"github.com/minipkg/log"
	"github.com/yaruz/app/pkg/yarus_platform/config"

	minipkg_gorm "github.com/minipkg/db/gorm"
)

type Sharding struct {
	IsAutoMigrate                bool
	Model                        interface{}
	Default                      *DBCluster
	BySysnames                   map[string]*DBCluster
	ClusterSysnamesByEntityTypes map[string]string
}

func newSharding(ctx context.Context, logger log.Logger, cfg *config.Sharding, model interface{}) (*Sharding, error) {
	clusterSysnamesByEntityTypes := make(map[string]string)
	defaultCluster, err := newDBCluster(logger, &cfg.Default)
	if err != nil {
		return nil, err
	}

	bySysnames := make(map[string]*DBCluster, len(cfg.BySysnames))
	for sysname, clusterCfg := range cfg.BySysnames {
		for _, entityType := range clusterCfg.EntityTypes {
			clusterSysnamesByEntityTypes[entityType] = sysname
		}

		cluster, err := newDBCluster(logger, &clusterCfg)
		if err != nil {
			return nil, err
		}
		bySysnames[sysname] = cluster
	}

	s := &Sharding{
		IsAutoMigrate:                cfg.IsAutoMigrate,
		Model:                        model,
		Default:                      defaultCluster,
		BySysnames:                   bySysnames,
		ClusterSysnamesByEntityTypes: clusterSysnamesByEntityTypes,
	}

	if err := s.SchemesInitWithContext(ctx, model); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Sharding) GetClustersByEntityTypes(entityTypes []string) ([]*DBCluster, error) {
	var DBClusters []*DBCluster

	if entityTypes == nil || len(entityTypes) == 0 {
		DBClusters = append(DBClusters, s.Default)
		return DBClusters, nil
	}

	for _, entityType := range entityTypes {
		cluster := s.Default
		if sysname, ok := s.ClusterSysnamesByEntityTypes[entityType]; ok {
			if cluster, ok = s.BySysnames[sysname]; !ok {
				return nil, errors.Wrapf(yaruserror.ErrNotFound, "Cluster with sysname = %q not found.", sysname)
			}
		}
		DBClusters = append(DBClusters, cluster)
	}

	return DBClusters, nil
}

func (s *Sharding) SchemesInitWithContext(ctx context.Context, model interface{}) (err error) {
	return s.ApplyFunc2DBs(func(db minipkg_gorm.IDB) (err error) {
		db, err = db.SchemeInitWithContext(ctx, model)
		return err
	})
}

func (s *Sharding) Close() (err error) {
	return s.ApplyFunc2DBs(func(db minipkg_gorm.IDB) error {
		return s.ApplyFunc2DBs(func(db minipkg_gorm.IDB) error {
			return db.Close()
		})
	})
}

func (s *Sharding) ApplyFunc2DBs(f func(db minipkg_gorm.IDB) error) (err error) {
	for _, shards := range s.BySysnames {
		if err = shards.ApplyFunc2DBs(f); err != nil {
			return err
		}
	}
	return s.Default.ApplyFunc2DBs(f)
}

type DBCluster struct {
	Capacity    uint
	EntityTypes []string
	Items       []minipkg_gorm.IDB
}

func newDBCluster(logger log.Logger, cfg *config.DBCluster) (*DBCluster, error) {
	var err error

	items := make([]minipkg_gorm.IDB, len(cfg.Items))
	for i, cfgItem := range cfg.Items {
		if items[i], err = minipkg_gorm.New(logger, cfgItem); err != nil {
			return nil, err
		}
	}

	return &DBCluster{
		Capacity:    cfg.Capacity,
		EntityTypes: cfg.EntityTypes,
		Items:       items,
	}, nil
}

func (s *DBCluster) SchemesInitWithContext(ctx context.Context, model interface{}) (err error) {
	return s.ApplyFunc2DBs(func(db minipkg_gorm.IDB) (err error) {
		db, err = db.SchemeInitWithContext(ctx, model)
		return err
	})
}

func (s *DBCluster) Close() (err error) {
	return s.ApplyFunc2DBs(func(db minipkg_gorm.IDB) error {
		return db.Close()
	})
}

func (s *DBCluster) ApplyFunc2DBs(f func(db minipkg_gorm.IDB) error) (err error) {
	for i := range s.Items {
		if err = f(s.Items[i]); err != nil {
			return err
		}
	}
	return nil
}
