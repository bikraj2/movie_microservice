package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"net/http"

	movie "bikraj.movie_microservice.net/movie/internal/controller"
	metadatagateway "bikraj.movie_microservice.net/movie/internal/gateway/metdata/http"
	ratinggateway "bikraj.movie_microservice.net/movie/internal/gateway/rating/http"
	httphandler "bikraj.movie_microservice.net/movie/internal/handler/http"
	"bikraj.movie_microservice.net/pkg/discovery"
	consul "bikraj.movie_microservice.net/pkg/discovery/consul"
)

const (
	serviceName = "movie"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Server listening on port %d", port)

	ctx := context.Background()
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("locahost:%d", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to Report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.DeRegister(ctx, instanceID, serviceName)

	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)

	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
