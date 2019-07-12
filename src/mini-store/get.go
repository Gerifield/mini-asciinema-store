package ministore

import (
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id",http.StatusBadRequest)
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid id",http.StatusBadRequest)
		return
	}

	f, err := os.Open(s.uploadPath+id)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, "file write error", http.StatusInternalServerError)
		return
	}
}

