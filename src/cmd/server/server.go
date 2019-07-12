package main

import (
	"flag"
	"github.com/gerifield/mini-asciinema-store/src/mini-store"
	"log"
	"net/http"
	"os"
)


func main()  {
	listenAddr := flag.String("listenAddr", ":8080", "HTTP listening address")

	uploadPath := flag.String("uploadPath", "uploads/", "Folder to store the uploaded files")
	serverBaseURL := flag.String("serverBaseURL", "http://127.0.0.1:8080", "Base URL for the server")

	https := flag.Bool("https", false, "HTTPS enable")
	httpsCert := flag.String("httpsCert", "server.crt", "HTTPS cert")
	httpsPrivateKey := flag.String("httpsPrivateKey", "server.key", "HTTPS private key")
	flag.Parse()

	srv := ministore.New(*uploadPath, *serverBaseURL)

	var err error
	if *https {
		err = http.ListenAndServeTLS(*listenAddr, *httpsCert, *httpsPrivateKey, srv.Routes())
	} else {
		err = http.ListenAndServe(*listenAddr, srv.Routes())
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

