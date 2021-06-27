package v1

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(db *sql.DB) (*chi.Mux, error) {
	bookHandler, err := NewBookHandler(db)
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
