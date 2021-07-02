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
	"go.uber.org/zap"
)

func authorRoutes(h *authorHandler) chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(
			cors.Handler(cors.Options{AllowedMethods: []string{"GET"}}),
		)
		r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			_ = render.Render(w, r, ErrMethodNotAllowed(r.Method))
			return
		})
		r.Get("/", h.List)
	})
	return r
}

/**********************************************************
 * Request and Response payloads/models for the REST api.
 **********************************************************/

// AuthorResponse is the response struct sent back to the client.
// Currently it embeds a pointer to entity.Author. In the future it would be
// possible to separate the two models and perform mapping if necessary.
type AuthorResponse struct {
	entity.Author
}

// newBookResponse accepts a pointer to an entity.Book and returns it embedded
// into a BookResponse.
func newAuthorResponse(author entity.Author) *AuthorResponse {
	return &AuthorResponse{Author: author}
}

// Render is a stub for preprocessing the BookResponse model. In the future it may
// be necessary to add some further data handling in here.
func (br *AuthorResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func newAuthorListResponse(authors []entity.Author) []render.Renderer {
	out := make([]render.Renderer, len(authors))
	for i, a := range authors {
		out[i] = newAuthorResponse(a)
	}
	return out
}

/*****************************
 * v1 Author endpoint handlers
 *****************************/

// authorHandler is used to wrap an author.Driver. This can be used by the API to
// drive the usecases powered by the driver.
type authorHandler struct {
	driver author.Driver
	logger *zap.Logger
}

// NewAuthorHandler accepts an author.Driver and, if valid, returns a pointer to an authorHandler. If the author.Driver
// is nil, then an error is returned.
func NewAuthorHandler(driver author.Driver, logger *zap.Logger) (*authorHandler, error) {
	if driver == nil {
		return nil, errors.New("non-nil author driver is required to create an author handler")
	}

	return &authorHandler{driver: driver, logger: logger}, nil
}

// List is an HTTP method that lists all entity.Author types that are accessible in the authorHandler's driver.
func (handler *authorHandler) List(w http.ResponseWriter, r *http.Request) {
	authors, err := handler.driver.ListAuthors(r.Context())
	if err != nil {
		handler.logger.Error(fmt.Sprintf("error listing authors: %s", err))
		_ = render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err := render.RenderList(w, r, newAuthorListResponse(authors)); err != nil {
		handler.logger.Error(fmt.Sprintf("error rendering authors: %s", err))
		_ = render.Render(w, r, ErrInternalServer(err))
		return
	}
}
