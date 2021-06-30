package v1

import (
	"fmt"
	"net/http"

	"github.com/LeviMatus/readcommend/service/internal/driver/genre"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

func genreRoutes(h *genreHandler) chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(
			cors.Handler(cors.Options{AllowedMethods: []string{"GET"}}),
		)
		r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, fmt.Sprintf("HTTP method %s is not allowed", r.Method), 400)
		})
		r.Get("/", h.List)
	})
	return r
}

// genreHandler is used to wrap an genre.Driver. This can be used by the API to
// drive the usecases powered by the driver.
type genreHandler struct {
	driver genre.Driver
}

// NewGenreHandler accepts an genre.Driver and, if valid, returns a pointer to an genreHandler. If the genre.Driver
// is nil, then an error is returned.
func NewGenreHandler(driver genre.Driver) (*genreHandler, error) {
	if driver == nil {
		return nil, errors.New("non-nil genre driver is required to create a genre handler")
	}

	return &genreHandler{driver: driver}, nil
}

// List is an HTTP method that lists all entity.Genre types that are accessible in the genreHandler's driver.
func (handler *genreHandler) List(w http.ResponseWriter, r *http.Request) {
	genres, err := handler.driver.ListGenres(r.Context())
	if err != nil {
		http.Error(w, "internal server error", 400)
		return
	}

	// An adaptor between the service layer and persistence layer
	// wouldn't be out of the question, but the conversion is very simple
	// so I'll just do it directly here. In the future, abstracting this
	// may be appropriate.
	var out = make([]entity.Genre, len(genres))
	for i, g := range genres {
		out[i] = entity.Genre{
			ID:    g.ID,
			Title: g.Title,
		}
	}

	render.JSON(w, r, out)
}
