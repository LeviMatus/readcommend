package v1

import (
	"net/http"

	"github.com/LeviMatus/readcommend/service/internal/driver/size"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

func sizeRoutes(h *sizeHandler) chi.Router {
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

// sizeHandler is used to wrap an size.Driver. This can be used by the API to
// drive the usecases powered by the driver.
type sizeHandler struct {
	driver size.Driver
}

// NewSizeHandler accepts an size.Driver and, if valid, returns a pointer to an sizeHandler. If the size.Driver
// is nil, then an error is returned.
func NewSizeHandler(driver size.Driver) (*sizeHandler, error) {
	if driver == nil {
		return nil, errors.New("non-nil size driver is required to create a size handler")
	}

	return &sizeHandler{driver: driver}, nil
}

// List is an HTTP method that lists all entity.Size types that are accessible in the sizeHandler's driver.
func (handler *sizeHandler) List(w http.ResponseWriter, r *http.Request) {
	sizes, err := handler.driver.ListSizes(r.Context())
	if err != nil {
		_ = render.Render(w, r, ErrInternalServer(err))
		return
	}

	// An adaptor between the service layer and persistence layer
	// wouldn't be out of the question, but the conversion is very simple
	// so I'll just do it directly here. In the future, abstracting this
	// may be appropriate.
	var out = make([]entity.Size, len(sizes))
	for i, e := range sizes {
		out[i] = entity.Size{
			ID:       e.ID,
			Title:    e.Title,
			MinPages: e.MinPages,
			MaxPages: e.MaxPages,
		}
	}

	render.JSON(w, r, out)
}
