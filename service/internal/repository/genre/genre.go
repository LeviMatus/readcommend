package genre

import (
	"context"
)

// Genre represents a category of a literary work.
type Genre struct {
	_ struct{}
	// ID is the Genre's primary key.
	ID    int
	Title string
}

// Repository defines the method set required for all Genre data-sources.
type Repository interface {
	// GetGenres returns all Genre entities in the DB resource.
	GetGenres(context.Context) ([]Genre, error)

	// Close should kill the DB connection.
	Close() error
}
