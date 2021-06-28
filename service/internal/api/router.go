package api

import (
	"fmt"
	"net/http"
	"time"

	v1 "github.com/LeviMatus/readcommend/service/internal/api/v1"
	"github.com/LeviMatus/readcommend/service/internal/driver"
	"github.com/LeviMatus/readcommend/service/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Server struct {
	mux *chi.Mux

	host string
	port string
}

func New(driver driver.Driver, config config.API) (*Server, error) {
	s := Server{
		mux:  chi.NewRouter(),
		host: config.Interface,
		port: config.Port,
	}

	s.mux.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(20*time.Second),
		render.SetContentType(render.ContentTypeJSON),
	)

	v1Router, err := v1.NewRouter(driver, config)
	if err != nil {
		return nil, err
	}

	s.mux.Route("/api", func(r chi.Router) {
		r.Mount("/v1", v1Router)
	})

	return &s, nil
}

func (s *Server) Listen() error {
	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.host, s.port), s.mux)
}
