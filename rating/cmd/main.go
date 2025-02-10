package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"bikraj.movie_microservice.net/gen"
	"bikraj.movie_microservice.net/pkg/discovery"
	"bikraj.movie_microservice.net/pkg/discovery/consul"
	"bikraj.movie_microservice.net/rating/internal/controller/rating"
	grpcHandler "bikraj.movie_microservice.net/rating/internal/handler/grpc"
	"bikraj.movie_microservice.net/rating/internal/repository/memory"
	"google.golang.org/grpc"
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
	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", "locahost:8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err.Error())
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	srv.Serve(lis)
	log.Println("Started the server on port: ", port)

}
