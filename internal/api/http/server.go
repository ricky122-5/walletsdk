package http

import (
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router *chi.Mux
}

func NewServer() *Server {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	return &Server{router: router}
}

func (s *Server) Router() *chi.Mux {
	return s.router
}

func (s *Server) ServeHTTP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	s.router.ServeHTTP(w, r)
}
