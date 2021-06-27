package book

import (
	"context"
)

// Book represents a literary work in a table/collection/node.
// This is a repository layer DTO.
type Book struct {
	_ struct{}
	// ID is the Book's primary key.
	ID            int32
	Title         string
	YearPublished int16
	Rating        float32
	Pages         int16
	GenreID       int16
	AuthorID      int16
}

type GetBooksParams struct {
	_ struct{}

	Title            *string
	MaxYearPublished *int16
	MinYearPublished *int16
	MaxPages         *int16
	MinPages         *int16
	Rating           *float32
	GenreIDs         []int16
	AuthorIDs        []int16
	Limit            *uint64
}

// Repository defines the method set required for all Book data-sources.
type Repository interface {
	// GetBooks returns all Book entities in the DB resource.
	GetBooks(context.Context, GetBooksParams) ([]Book, error)

	// Close should kill the DB connection.
	Close() error
}
