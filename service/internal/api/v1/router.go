package v1

import (
	"fmt"
	"net/http"

	"github.com/LeviMatus/readcommend/service/internal/driver"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	return r, nil
}
