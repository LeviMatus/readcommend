package size

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/encoding"
)

// Size represents a category of a literary work.
type Size struct {
	// ID is the Size's primary key.
	ID int

	// Title represents the size category.
	Title string

	// MinimumPages is the lower-bound of the allowed pages for this Size category.
	MinimumPages encoding.NullInt16

	// MaximumPages is the upper-bound of the allowed pages for this Size category.
	MaximumPages encoding.NullInt16
}

// Repository defines the method set required for all Size data-sources.
type Repository interface {
	// GetSizes returns all Size entities in the DB resource.
	GetSizes(context.Context) ([]Size, error)

	// Close should kill the DB connection.
	Close() error
}
