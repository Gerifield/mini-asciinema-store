package ministore

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := uuid.Parse(id)
	if err != nil {
		respondErr(err, w, "invalid id", http.StatusBadRequest)
		return
	}

	br, err := s.uploadBucket.NewReader(r.Context(), id, nil)
	if err != nil {
		respondErr(err, w, "file not found", http.StatusNotFound)
		return
	}
	defer br.Close()

	_, err = io.Copy(w, br)
	if err != nil {
		respondErr(err, w, "file write error", http.StatusInternalServerError)
		return
	}
}
