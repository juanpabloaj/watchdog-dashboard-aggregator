package main

import (
	"log"
	"net/http"
)

var (
	buildDate string
	gitHash   string
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	port := ":8080"

	log.Printf("build date: %s", buildDate)
	log.Printf("git hash: %s", gitHash)

	log.Printf("starting server on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
