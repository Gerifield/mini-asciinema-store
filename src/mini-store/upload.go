package ministore

import (
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	id := uuid.New()
	var hasWriteError bool

	bw, err := s.uploadBucket.NewWriter(r.Context(), id.String(), nil)
	if err != nil {
		respondErr(err, w, "file create failed", http.StatusInternalServerError)
		return
	}
	defer func() {
		err := bw.Close()
		if err != nil && !hasWriteError {
			respondErr(err, w, "file write commit failed", http.StatusInternalServerError)
		}
	}()

	mf, _, err := r.FormFile("asciicast")
	if err != nil {
		hasWriteError = true
		respondErr(err, w, "file read failed", http.StatusInternalServerError)
		return
	}
	defer mf.Close()

	//fmt.Println("size:",mfh.Size, "name:", mfh.Filename)

	_, err = io.Copy(bw, mf)
	if err != nil {
		log.Println(err)
		hasWriteError = true
		respondErr(err, w, "file save failed", http.StatusInternalServerError)
		return
	}

	//fmt.Println(string(b))
	w.Write([]byte(s.baseURL + "/a/" + id.String()))
}
