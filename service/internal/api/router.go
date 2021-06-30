package api

import (
	"net"
	"net/http"

	v1 "github.com/LeviMatus/readcommend/service/internal/api/v1"
	"github.com/LeviMatus/readcommend/service/internal/driver"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

type Server struct {
	mux *chi.Mux

	host string
	port string
}

func New(driver driver.Driver) (*Server, error) {
	if driver == nil {
		return nil, errors.New("a non-nil driver is required")
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

	v1Router, err := v1.NewRouter(driver)
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
