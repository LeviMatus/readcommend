package era

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/encoding"
)

// Era represents a literary era in a table/collection/node.
// This is a repository layer DTO.
type Era struct {
	_ struct{}
	// ID is the Era's primary key.
	ID int

	// Title represents the name of the Era.
	Title string

	// StartYear is the year in which the literary Era begins.
	StartYear encoding.NullInt16

	// EndYear is the year in which the literary Era ends.
	EndYear encoding.NullInt16
}

// Repository defines the method set required for all Era data-sources.
type Repository interface {
	// GetEras returns all Era entities in the DB resource.
	GetEras(context.Context) ([]Era, error)

	// Close should kill the DB connection.
	Close() error
}
