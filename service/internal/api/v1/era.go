package v1

import (
	"net/http"

	"github.com/LeviMatus/readcommend/service/internal/driver/era"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

func eraRoutes(h *eraHandler) chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(
			cors.Handler(cors.Options{AllowedMethods: []string{"GET"}}),
		)
		r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			_ = render.Render(w, r, ErrMethodNotAllowed(r.Method))
		})
		r.Get("/", h.List)
	})
	return r
}

// eraHandler is used to wrap an era.Driver. This can be used by the API to
// drive the usecases powered by the driver.
type eraHandler struct {
	driver era.Driver
}

// NewEraHandler accepts an era.Driver and, if valid, returns a pointer to an eraHandler. If the era.Driver
// is nil, then an error is returned.
func NewEraHandler(driver era.Driver) (*eraHandler, error) {
	if driver == nil {
		return nil, errors.New("non-nil era driver is required to create a era handler")
	}

	return &eraHandler{driver: driver}, nil
}

// List is an HTTP method that lists all entity.Era types that are accessible in the eraHandler's driver.
func (handler *eraHandler) List(w http.ResponseWriter, r *http.Request) {
	eras, err := handler.driver.ListEras(r.Context())
	if err != nil {
		_ = render.Render(w, r, ErrInternalServer(err))
		return
	}

	// An adaptor between the service layer and persistence layer
	// wouldn't be out of the question, but the conversion is very simple
	// so I'll just do it directly here. In the future, abstracting this
	// may be appropriate.
	var out = make([]entity.Era, len(eras))
	for i, e := range eras {
		out[i] = entity.Era{
			ID:      e.ID,
			Title:   e.Title,
			MinYear: e.MinYear,
			MaxYear: e.MaxYear,
		}
	}

	render.JSON(w, r, out)
}
