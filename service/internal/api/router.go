package api

import (
	"net"
	"net/http"

	v1 "github.com/LeviMatus/readcommend/service/internal/api/v1"
	"github.com/LeviMatus/readcommend/service/internal/driver/author"
	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/driver/era"
	"github.com/LeviMatus/readcommend/service/internal/driver/genre"
	"github.com/LeviMatus/readcommend/service/internal/driver/size"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Server struct {
	mux *chi.Mux

	host string
	port string
}

type RequiredDrivers struct {
	AuthorDriver *author.Driver
	GenreDriver  *genre.Driver
	EraDriver    *era.Driver
	SizeDriver   *size.Driver
	BookDriver   *book.Driver
}

// Validate ensures that the required drivers are provided.
func (req RequiredDrivers) Validate() error {
	if req.AuthorDriver == nil ||
		req.SizeDriver == nil ||
		req.GenreDriver == nil ||
		req.EraDriver == nil ||
		req.BookDriver == nil {
		return errors.New("all drivers must be non-nil")
	}
	return nil
}

func New(ad author.Driver, sd size.Driver, gd genre.Driver, ed era.Driver, bd book.Driver, logger *zap.Logger) (*Server, error) {
	if ad == nil || sd == nil || gd == nil || ed == nil || bd == nil || logger == nil {
		return nil, errors.New("dependencies for the API are not satisfied - non-nil drivers and logger are required")
	}

	s := Server{
		mux: chi.NewRouter(),
	}

	s.mux.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		render.SetContentType(render.ContentTypeJSON),
	)

	v1Router, err := v1.NewRouter(ad, sd, gd, ed, bd, logger)
	if err != nil {
		return nil, err
	}

	s.mux.Route("/api", func(r chi.Router) {
		r.Mount("/v1", v1Router)
	})

	return &s, nil
}

func (s *Server) Serve(listener net.Listener) error {
	return http.Serve(listener, s.mux)
}
