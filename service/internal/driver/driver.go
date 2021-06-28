package driver

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/driver/author"
	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/driver/era"
	"github.com/LeviMatus/readcommend/service/internal/driver/genre"
	"github.com/LeviMatus/readcommend/service/internal/driver/size"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/pkg/errors"
)

type Driver interface {
	author.Driver
	genre.Driver
	size.Driver
	era.Driver
	book.Driver
}

type driver struct {
	authorDriver author.Driver
	genreDriver  genre.Driver
	sizeDriver   size.Driver
	eraDriver    era.Driver
	bookDriver   book.Driver
}

func New(a author.Driver, g genre.Driver, s size.Driver, e era.Driver, b book.Driver) (*driver, error) {
	if a == nil || g == nil || s == nil || e == nil || b == nil {
		return nil, errors.New("")
	}

	return &driver{
		authorDriver: a,
		genreDriver:  g,
		sizeDriver:   s,
		eraDriver:    e,
		bookDriver:   b,
	}, nil
}

func (d driver) ListAuthors(ctx context.Context) ([]entity.Author, error) {
	return d.authorDriver.ListAuthors(ctx)
}

func (d driver) ListGenres(ctx context.Context) ([]entity.Genre, error) {
	return d.genreDriver.ListGenres(ctx)
}

func (d driver) ListSizes(ctx context.Context) ([]entity.Size, error) {
	return d.sizeDriver.ListSizes(ctx)
}

func (d driver) ListEras(ctx context.Context) ([]entity.Era, error) {
	return d.eraDriver.ListEras(ctx)
}

func (d driver) SearchBooks(ctx context.Context, params book.SearchInput) ([]entity.Book, error) {
	return d.bookDriver.SearchBooks(ctx, params)
}
