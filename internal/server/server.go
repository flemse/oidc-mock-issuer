package server

import (
	"log"
	"net/http"
)

// Start starts the HTTP server on the given address.
func Start(addr string, handler http.Handler) {
	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
