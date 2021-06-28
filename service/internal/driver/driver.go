package driver

import (
	"github.com/LeviMatus/readcommend/service/internal/driver/author"
	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/driver/era"
	"github.com/LeviMatus/readcommend/service/internal/driver/genre"
	"github.com/LeviMatus/readcommend/service/internal/driver/size"
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
