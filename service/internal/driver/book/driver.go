package book

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
)

type SearchInput struct {
	_ struct{}

	Title            *string
	MaxYearPublished *int16
	MinYearPublished *int16
	MaxPages         *int16
	MinPages         *int16
	GenreIDs         []int16
	AuthorIDs        []int16
	Limit            *uint64
}

type driver struct {
	repository Repository
}

func NewDriver(r Repository) *driver {
	return &driver{repository: r}
}

func (d *driver) Search(ctx context.Context, params SearchInput) ([]entity.Book, error) {
	return d.repository.Search(ctx, params)
}
