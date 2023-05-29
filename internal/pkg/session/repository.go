package session

import (
	"context"
)

// Repository encapsulates the logic to access session from the data source.
type Repository interface {
	// Get returns the session with the specified user ID.
	Get(ctx context.Context, ID string) (*Session, error)
	// Create saves a new entity in the storage.
	Set(ctx context.Context, entity *Session) error
	// Delete removes the entity with given ID from the storage.
	Delete(ctx context.Context, ID string) error
}
