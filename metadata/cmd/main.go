package main

import (
	"log"
	"net/http"

	"bikraj.movie_microservice.net/metadata/internal/controller/metadata"
	httphandler "bikraj.movie_microservice.net/metadata/internal/handler"
	"bikraj.movie_microservice.net/metadata/internal/repository/memory"
)

func main() {
	log.Println("Starting the movie metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}
