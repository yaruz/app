package entity

import (
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"
)

type RelationValue struct {
	entity_type.Relation
	Value []uint
}
