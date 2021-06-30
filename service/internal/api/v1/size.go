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

/**********************************************************
 * Request and Response payloads/models for the REST api.
 **********************************************************/

// SizeResponse is the response struct sent back to the client.
// Currently it embeds a pointer to entity.Size. In the future it would be
// possible to separate the two models and perform mapping if necessary.
type SizeResponse struct {
	entity.Size
}

// newBookResponse accepts a pointer to an entity.Book and returns it embedded
// into a BookResponse.
func newSizeResponse(size entity.Size) *SizeResponse {
	return &SizeResponse{Size: size}
}

// Render is a stub for preprocessing the BookResponse model. In the future it may
// be necessary to add some further data handling in here.
func (br *SizeResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func newSizeListResponse(sizes []entity.Size) []render.Renderer {
	out := make([]render.Renderer, len(sizes))
	for i, s := range sizes {
		out[i] = newSizeResponse(s)
	}
	return out
}

/*****************************
 * v1 Size endpoint handlers
 *****************************/

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

	if err := render.RenderList(w, r, newSizeListResponse(sizes)); err != nil {
		_ = render.Render(w, r, ErrInternalServer(err))
		return
	}
}
