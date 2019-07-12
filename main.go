package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

const uploadPath = "uploads/"
const serverBaseURL = "http://127.0.0.1:8080"

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	id := uuid.New()

	f, err := os.Create(uploadPath+id.String())
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
	w.Write([]byte(serverBaseURL+"/a/"+id.String()))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
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

	f, err := os.Open(uploadPath+id)
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

func main()  {
	listenAddr := flag.String("listenAddr", ":8080", "HTTP listening address")
	flag.Parse()

	mux := chi.NewRouter()
	mux.Post("/api/asciicasts", uploadHandler)
	mux.Get("/a/{id}", getHandler)

	err := http.ListenAndServe(*listenAddr, mux)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
