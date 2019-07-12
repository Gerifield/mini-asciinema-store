package ministore

import (
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
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

func (s *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	id := uuid.New()

	f, err := os.Create(s.uploadPath+id.String())
	if err != nil {
		http.Error(w, "file create failed", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	mf, _, err := r.FormFile("asciicast")
	if err != nil {
		http.Error(w, "file read failed", http.StatusInternalServerError)
		return
	}
	defer mf.Close()

	//fmt.Println("size:",mfh.Size, "name:", mfh.Filename)

	_, err = io.Copy(f, mf)
	if err != nil {
		http.Error(w, "file save failed", http.StatusInternalServerError)
		return
	}

	//fmt.Println(string(b))
	w.Write([]byte(s.baseURL+"/a/"+id.String()))
}

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
