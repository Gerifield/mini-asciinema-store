package ministore

import (
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
)

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