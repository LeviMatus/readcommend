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

/**********************************************************
 * Request and Response payloads/models for the REST api.
 **********************************************************/

// EraResponse is the response struct sent back to the client.
// Currently it embeds a pointer to entity.Era. In the future it would be
// possible to separate the two models and perform mapping if necessary.
type EraResponse struct {
	entity.Era
}

// newBookResponse accepts a pointer to an entity.Book and returns it embedded
// into a BookResponse.
func newEraResponse(era entity.Era) *EraResponse {
	return &EraResponse{Era: era}
}

// Render is a stub for preprocessing the BookResponse model. In the future it may
// be necessary to add some further data handling in here.
func (br *EraResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func newEraListResponse(eras []entity.Era) []render.Renderer {
	out := make([]render.Renderer, len(eras))
	for i, e := range eras {
		out[i] = newEraResponse(e)
	}
	return out
}

/*****************************
 * v1 Era endpoint handlers
 *****************************/

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

	if err := render.RenderList(w, r, newEraListResponse(eras)); err != nil {
		_ = render.Render(w, r, ErrInternalServer(err))
		return
	}
}
