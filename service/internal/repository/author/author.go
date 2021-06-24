package author

import (
	"context"
)

// Author represents a literary writer in a table/collection/node.
// This is a repository layer DTO.
type Author struct {
	// ID is the Author's primary key.
	ID        int
	FirstName string
	LastName  string
}

// Repository defines the method set required for all Author data-sources.
type Repository interface {
	// GetAuthors returns all Author entities in the DB resource.
	GetAuthors(context.Context) ([]Author, error)

	// Close should kill the DB connection.
	Close() error
}
