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
	"go.uber.org/zap"
)

func genreRoutes(h *genreHandler) chi.Router {
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

// GenreResponse is the response struct sent back to the client.
// Currently it embeds a pointer to entity.Genre. In the future it would be
// possible to separate the two models and perform mapping if necessary.
type GenreResponse struct {
	entity.Genre
}

// newBookResponse accepts a pointer to an entity.Book and returns it embedded
// into a BookResponse.
func newGenreResponse(genre entity.Genre) *GenreResponse {
	return &GenreResponse{Genre: genre}
}

// Render is a stub for preprocessing the BookResponse model. In the future it may
// be necessary to add some further data handling in here.
func (br *GenreResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func newGenreListResponse(genres []entity.Genre) []render.Renderer {
	out := make([]render.Renderer, len(genres))
	for i, g := range genres {
		out[i] = newGenreResponse(g)
	}
	return out
}

/*****************************
 * v1 Genre endpoint handlers
 *****************************/

// genreHandler is used to wrap an genre.Driver. This can be used by the API to
// drive the usecases powered by the driver.
type genreHandler struct {
	driver genre.Driver
	logger *zap.Logger
}

// NewGenreHandler accepts an genre.Driver and, if valid, returns a pointer to an genreHandler. If the genre.Driver
// is nil, then an error is returned.
func NewGenreHandler(driver genre.Driver, logger *zap.Logger) (*genreHandler, error) {
	if driver == nil {
		return nil, errors.New("non-nil genre driver is required to create a genre handler")
	}

	return &genreHandler{driver: driver, logger: logger}, nil
}

// List is an HTTP method that lists all entity.Genre types that are accessible in the genreHandler's driver.
func (handler *genreHandler) List(w http.ResponseWriter, r *http.Request) {
	genres, err := handler.driver.ListGenres(r.Context())
	if err != nil {
		handler.logger.Error(fmt.Sprintf("error listing genres: %s", err))
		_ = render.Render(w, r, ErrInternalServer(err))
		return
	}

	if err := render.RenderList(w, r, newGenreListResponse(genres)); err != nil {
		handler.logger.Error(fmt.Sprintf("error rendering eras: %s", err))
		_ = render.Render(w, r, ErrInternalServer(err))
		return
	}
}
