package session

import (
	"context"
)

// Repository encapsulates the logic to access session from the data source.
type Repository interface {
	// GetByUserID returns the session with the specified user ID.
	Get(ctx context.Context, userId uint) (*Session, error)
	// Create saves a new entity in the storage.
	Create(ctx context.Context, entity *Session) error
	// Update updates the entity with given ID in the storage.
	Update(ctx context.Context, entity *Session) error
	// Delete removes the entity with given ID from the storage.
	Delete(ctx context.Context, entity *Session) error
}
