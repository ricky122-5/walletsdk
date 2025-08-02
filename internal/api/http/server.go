package http

import (
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server wraps routing for the wallet SDK API.
type Server struct {
	router *chi.Mux
}

// NewServer constructs the HTTP server with standard middleware.
func NewServer() *Server {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	return &Server{router: router}
}

// Router exposes the underlying chi router for mounting handlers.
func (s *Server) Router() *chi.Mux {
	return s.router
}

// ServeHTTP allows the server to satisfy http.Handler.
func (s *Server) ServeHTTP(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	s.router.ServeHTTP(w, r)
}
