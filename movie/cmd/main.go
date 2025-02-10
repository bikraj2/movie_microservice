package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"bikraj.movie_microservice.net/gen"
	movie "bikraj.movie_microservice.net/movie/internal/controller"
	metadatagateway "bikraj.movie_microservice.net/movie/internal/gateway/metdata/grpc"
	ratinggateway "bikraj.movie_microservice.net/movie/internal/gateway/rating/grpc"
	grpchandler "bikraj.movie_microservice.net/movie/internal/handler/grpc"
	"bikraj.movie_microservice.net/pkg/discovery"
	consul "bikraj.movie_microservice.net/pkg/discovery/consul"
	"google.golang.org/grpc"
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

	svc := movie.New(ratingGateway, metadataGateway)
	h := grpchandler.New(svc)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("err while listening: %v", err.Error())
	}
	srv := grpc.NewServer()
	gen.RegisterMovieServiceServer(srv, h)
	srv.Serve(lis)
}
