package ministore

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"gocloud.dev/blob"
)

// Server .
type Server struct {
	baseURL      string
	uploadBucket *blob.Bucket
}

// New mini-store server init
func New(baseURL string, uploadBucket *blob.Bucket) *Server {
	return &Server{
		baseURL:      baseURL,
		uploadBucket: uploadBucket,
	}
}

// Routes returns the router
func (s *Server) Routes() http.Handler {
	mux := chi.NewRouter()
	mux.Post("/api/asciicasts", s.uploadHandler)
	mux.Get("/a/{id}", s.getHandler)
	return mux
}

func respondErr(err error, w http.ResponseWriter, errorStr string, code int) {
	http.Error(w, errorStr, code)
	log.Println(errorStr, err)
}
