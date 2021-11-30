package infrastructure

import (
	"context"

	"github.com/minipkg/selection_condition"

	minipkg_gorm "github.com/minipkg/db/gorm"
)

type Sharding struct {
	IsAutoMigrate bool
	Model         interface{}
	Default       Shards
	ByTypes       map[string]Shards
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
	for _, shards := range s.ByTypes {
		if err = shards.ApplyFunc2DBs(f); err != nil {
			return err
		}
	}
	return s.Default.ApplyFunc2DBs(f)
}

type Shards struct {
	Capacity uint
	Items    []minipkg_gorm.IDB
}

func (s *Shards) SchemesInitWithContext(ctx context.Context, model interface{}) (err error) {
	return s.ApplyFunc2DBs(func(db minipkg_gorm.IDB) (err error) {
		db, err = db.SchemeInitWithContext(ctx, model)
		return err
	})
}

func (s *Shards) Close() (err error) {
	return s.ApplyFunc2DBs(func(db minipkg_gorm.IDB) error {
		return db.Close()
	})
}

func (s *Shards) ApplyFunc2DBs(f func(db minipkg_gorm.IDB) error) (err error) {
	for i := range s.Items {
		if err = f(s.Items[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *Sharding) GetDBForInsert(typeID uint) minipkg_gorm.IDB {
	return nil
}

func (s *Sharding) GetDB(ID uint, typeID uint) minipkg_gorm.IDB {
	return nil
}

func (s *Sharding) GetDBs(condition *selection_condition.SelectionCondition) []minipkg_gorm.IDB {
	return nil
}
