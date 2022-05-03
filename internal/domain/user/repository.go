package user

import (
	"context"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	New(ctx context.Context) (*User, error)
	// Get returns the entity with the specified ID.
	Get(ctx context.Context, id uint, langID uint) (*User, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]User, error)
	First(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (*User, error)
	// Count returns the number of items.
	Count(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (uint, error)
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Create saves a new album in the storage.
	Create(ctx context.Context, obj *User, langID uint) error
	// Update updates the album with given ID in the storage.
	Update(ctx context.Context, obj *User, langID uint) error
	// Delete removes the album with given ID from the storage.
	Delete(ctx context.Context, id uint) error
	LangIDValidate(ctx context.Context, langID uint) error
}
