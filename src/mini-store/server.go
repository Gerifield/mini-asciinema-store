package ministore

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"gocloud.dev/blob"
)

// Server .
type Server struct {
	baseURL      string
	uploadBucket *blob.Bucket
	authFile     string
}

// New mini-store server init
func New(baseURL string, uploadBucket *blob.Bucket, authFile string) *Server {
	return &Server{
		baseURL:      baseURL,
		uploadBucket: uploadBucket,
		authFile:     authFile,
	}
}

// Routes returns the router
func (s *Server) Routes() http.Handler {
	mux := chi.NewRouter()

	if s.authFile != "" {
		tokens, err := readTokens(s.authFile)
		if err != nil {
			log.Fatalln("can't read the token file", err)
		} else {
			log.Println("loaded", len(tokens), "tokens")
		}
		mux.Use(authMiddleware(tokens))
	}
	mux.Use(simpleLogMiddleware)

	mux.Post("/api/asciicasts", s.uploadHandler)
	mux.Get("/a/{id}", s.getHandler)
	return mux
}

func readTokens(authFile string) ([]string, error) {
	b, err := ioutil.ReadFile(authFile)
	if err != nil {
		return nil, err
	}
	rawTokens := strings.Split(string(b), "\n")

	// cleanup empty lines
	tokens := make([]string, 0, len(rawTokens))
	var trimmed string
	for _, t := range rawTokens {
		trimmed = strings.TrimSpace(t)
		if len(trimmed) > 0 {
			tokens = append(tokens, trimmed)
		}
	}
	return tokens, nil
}

func authMiddleware(tokens []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, pass, ok := r.BasicAuth()
			if !ok {
				respondErr(nil, w, "invalid Authorization field", http.StatusUnauthorized)
				return
			}

			for _, t := range tokens {
				if t == pass {
					next.ServeHTTP(w, r)
					return
				}
			}

			respondErr(nil, w, "invalid token", http.StatusUnauthorized)
		})
	}
}

func simpleLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func respondErr(err error, w http.ResponseWriter, errorStr string, code int) {
	http.Error(w, errorStr, code)
	log.Println(errorStr, err)
}
