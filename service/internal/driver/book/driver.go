package book

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
)

// SearchInput is a input parameter for SearchBooks.
type SearchInput struct {
	_ struct{}

	// Title is used to search for books by title (ignored if nil).
	Title *string

	// MaxYearPublished filters Books published later than the specified year (ignored if nil).
	MaxYearPublished *int16

	// MinYearPublished filters Books published earlier than the specified year (ignored if nil).
	MinYearPublished *int16

	// MaxPages filters Books with more than the specified number of pages (ignored if nil).
	MaxPages *int16

	// MaxPages filters Books with less than the specified number of pages (ignored if nil).
	MinPages *int16

	// GenreIDs includes Books whose Genre's ID falls into _any_ of the included GenreIDs (ignored if nil).
	GenreIDs []int16

	// AuthorIDs includes Books whose Author's ID falls into _any_ of the included AuthorIDs (ignored if nil).
	AuthorIDs []int16

	// Limit specifies a maximum number of Books to be returned. If not specified, there will be no limit.
	Limit *uint64
}

type driver struct {
	repository Repository
}

// NewDriver creates a driver which wraps the repository. The wrapper
// will perform business logic against the usecases of Size entity.
func NewDriver(r Repository) *driver {
	return &driver{repository: r}
}

// SearchBooks searches for entity.Book types from the repository and returns them.
func (d *driver) SearchBooks(ctx context.Context, params SearchInput) ([]entity.Book, error) {
	return d.repository.Search(ctx, params)
}
