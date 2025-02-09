package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"bikraj.movie_microservice.net/metadata/internal/controller/metadata"
	httphandler "bikraj.movie_microservice.net/metadata/internal/handler"
	"bikraj.movie_microservice.net/metadata/internal/repository/memory"
	"bikraj.movie_microservice.net/pkg/discovery"
	"bikraj.movie_microservice.net/pkg/discovery/consul"
)

const servicename = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "api Handler port")
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
				log.Println("F    log.Printlniled to report healthy state: ", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, servicename)
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		panic(err)
	}
	log.Println("Started the server on port: ", port)
}
