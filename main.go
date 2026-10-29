package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	if err := InitCache(); err != nil {
		log.Fatalf("Could not initialize data: %v", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/artist", ArtistHandler)
	http.HandleFunc("/search", SearchHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
