package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	_ "gocloud.dev/blob/s3blob"

	"github.com/gerifield/mini-asciinema-store/src/mini-store"
)

func main() {
	listenAddr := flag.String("listenAddr", ":8080", "HTTP listening address")

	uploadBucket := flag.String("uploadBucket", "file:///uploads/", "Folder or bucket URL to store the uploaded files (supports: file, mem, s3)")
	serverBaseURL := flag.String("serverBaseURL", "http://127.0.0.1:8080", "Base URL for the server")

	https := flag.Bool("https", false, "HTTPS enable")
	httpsCert := flag.String("httpsCert", "server.crt", "HTTPS cert")
	httpsPrivateKey := flag.String("httpsPrivateKey", "server.key", "HTTPS private key")
	flag.Parse()

	bucket, err := blob.OpenBucket(context.Background(), *uploadBucket)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer bucket.Close()

	srv := ministore.New(*serverBaseURL, bucket)

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
