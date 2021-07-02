package v1

import (
	"fmt"

	"github.com/LeviMatus/readcommend/service/internal/driver/author"
	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/driver/era"
	"github.com/LeviMatus/readcommend/service/internal/driver/genre"
	"github.com/LeviMatus/readcommend/service/internal/driver/size"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(ad author.Driver, sd size.Driver, gd genre.Driver, ed era.Driver, bd book.Driver, logger *zap.Logger) (*chi.Mux, error) {
	bookHandler, err := NewBookHandler(bd, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to create v1 routes: %w", err)
	}

	authorHandler, err := NewAuthorHandler(ad, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to create v1 routes: %w", err)
	}

	genreHandler, err := NewGenreHandler(gd, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to create v1 routes: %w", err)
	}

	eraHandler, err := NewEraHandler(ed, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to create v1 routes: %w", err)
	}

	sizeHandler, err := NewSizeHandler(sd, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to create v1 routes: %w", err)
	}

	r := chi.NewRouter()

	r.Mount("/books", bookRoutes(bookHandler))
	r.Mount("/authors", authorRoutes(authorHandler))
	r.Mount("/genres", genreRoutes(genreHandler))
	r.Mount("/eras", eraRoutes(eraHandler))
	r.Mount("/sizes", sizeRoutes(sizeHandler))

	return r, nil
}
