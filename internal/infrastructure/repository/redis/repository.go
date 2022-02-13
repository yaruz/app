package redis

import (
	"github.com/minipkg/selection_condition"

	"github.com/minipkg/db/redis"
)

// IRepository is an interface of repository
type IRepository interface{}

// repository persists albums in database
type repository struct {
	db         redis.IDB
	Conditions selection_condition.SelectionCondition
}
