package v1

import (
	"fmt"
	"net/http"

	driver2 "github.com/LeviMatus/readcommend/service/internal/driver"
	"github.com/LeviMatus/readcommend/service/pkg/config"
	"github.com/go-chi/chi/v5"
)

func NewRouter(driver driver2.Driver, config config.API) (*chi.Mux, error) {
	bookHandler, err := NewBookHandler(driver, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create v1 routes: %w", err)
	}
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Mount("/books", func() http.Handler {
			br := chi.NewRouter()
			br.Get("/", bookHandler.List)
			return br
		}())
	})

	return r, nil
}
