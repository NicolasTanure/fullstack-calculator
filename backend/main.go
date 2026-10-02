package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"fullstack-calculator/backend/internal/httpapi"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           httpapi.NewHandler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Starting calculator API on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
