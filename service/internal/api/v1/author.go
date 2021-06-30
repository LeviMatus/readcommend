package v1

import (
	"fmt"
	"net/http"

	"github.com/LeviMatus/readcommend/service/internal/driver/author"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

func authorRoutes(h *authorHandler) chi.Router {
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

// authorHandler is used to wrap an author.Driver. This can be used by the API to
// drive the usecases powered by the driver.
type authorHandler struct {
	driver author.Driver
}

// NewAuthorHandler accepts an author.Driver and, if valid, returns a pointer to an authorHandler. If the author.Driver
// is nil, then an error is returned.
func NewAuthorHandler(driver author.Driver) (*authorHandler, error) {
	if driver == nil {
		return nil, errors.New("non-nil author driver is required to create an author handler")
	}

	return &authorHandler{driver: driver}, nil
}

// List is an HTTP method that lists all entity.Author types that are accessible in the authorHandler's driver.
func (a *authorHandler) List(w http.ResponseWriter, r *http.Request) {
	authors, err := a.driver.ListAuthors(r.Context())
	if err != nil {
		http.Error(w, "internal server error", 400)
		return
	}

	// An adaptor between the service layer and persistence layer
	// wouldn't be out of the question, but the conversion is very simple
	// so I'll just do it directly here. In the future, abstracting this
	// may be appropriate.
	var out = make([]entity.Author, len(authors))
	for i, b := range authors {
		out[i] = entity.Author{
			ID:        b.ID,
			FirstName: b.FirstName,
			LastName:  b.LastName,
		}
	}

	render.JSON(w, r, out)
}
