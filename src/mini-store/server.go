package ministore

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Server .
type Server struct {
	uploadPath string
	baseURL string
}

// New mini-store server init
func New(uploadPath string, baseURL string) *Server {
	return &Server{
		uploadPath: uploadPath,
		baseURL:    baseURL,
	}
}

// Routes returns the router
func (s *Server) Routes() http.Handler{
	mux := chi.NewRouter()
	mux.Post("/api/asciicasts", s.uploadHandler)
	mux.Get("/a/{id}", s.getHandler)
	return mux
}
