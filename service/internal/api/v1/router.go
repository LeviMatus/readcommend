package v1

import (
	"fmt"

	"github.com/LeviMatus/readcommend/service/internal/driver"
	"github.com/go-chi/chi/v5"
)

func NewRouter(driver driver.Driver) (*chi.Mux, error) {
	bookHandler, err := NewBookHandler(driver)
	if err != nil {
		return nil, fmt.Errorf("unable to create v1 routes: %w", err)
	}

	authorHandler, err := NewAuthorHandler(driver)
	if err != nil {
		return nil, fmt.Errorf("unable to create v1 routes: %w", err)
	}

	genreHandler, err := NewGenreHandler(driver)
	if err != nil {
		return nil, fmt.Errorf("unable to create v1 routes: %w", err)
	}

	sizeHandler, err := NewSizeHandler(driver)
	if err != nil {
		return nil, fmt.Errorf("unable to create v1 routes: %w", err)
	}

	r := chi.NewRouter()

	r.Mount("/books", func() http.Handler {
		br := chi.NewRouter()
		br.Use(
			cors.Handler(cors.Options{AllowedMethods: []string{"GET"}}),
			ValidateGetBookParams,
		)
		br.Get("/", bookHandler.List)
		return br
	}())

	r.Mount("/authors", authorRoutes(authorHandler))
	r.Mount("/genres", genreRoutes(genreHandler))
	r.Mount("/sizes", sizeRoutes(sizeHandler))

	return r, nil
}
