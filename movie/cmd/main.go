package main

import (
	"log"
	"net/http"

	movie "bikraj.movie_microservice.net/movie/internal/controller"
	metadatagateway "bikraj.movie_microservice.net/movie/internal/gateway/metdata/http"
	ratinggateway "bikraj.movie_microservice.net/movie/internal/gateway/rating/http"
	httphandler "bikraj.movie_microservice.net/movie/internal/handler/http"
)

func main() {
	log.Printf("Server listening on port 8083")

	metadataGateway := metadatagateway.New("localhost:8081")
	ratingGateway := ratinggateway.New("localhost:8082")

	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
