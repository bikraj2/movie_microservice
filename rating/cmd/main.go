package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"bikraj.movie_microservice.net/pkg/discovery"
	"bikraj.movie_microservice.net/pkg/discovery/consul"
	"bikraj.movie_microservice.net/rating/internal/controller/rating"
	httphandler "bikraj.movie_microservice.net/rating/internal/handler/http"
	"bikraj.movie_microservice.net/rating/internal/repository/memory"
)

const (
	servicename = "rating"
)

func main() {

	var port int
	flag.IntVar(&port, "port", 8082, "api Handler port")
	flag.Parse()
	log.Println("Starting the movie metadata service")
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	instanceID := discovery.GenerateInstanceID(servicename)

	if err := registry.Register(ctx, instanceID, servicename, fmt.Sprintf("locahhost:%v", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, servicename); err != nil {
				log.Println("Failed to report healthy state: ", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.DeRegister(ctx, instanceID, servicename)
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.Handle))
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		panic(err)
	}
	log.Println("Started the server on port: ", port)
}
