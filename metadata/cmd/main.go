package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"bikraj.movie_microservice.net/gen"
	"bikraj.movie_microservice.net/metadata/internal/controller/metadata"
	grpchandler "bikraj.movie_microservice.net/metadata/internal/handler/grpc"
	"bikraj.movie_microservice.net/metadata/internal/repository/memory"
	"bikraj.movie_microservice.net/pkg/discovery"
	"bikraj.movie_microservice.net/pkg/discovery/consul"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

const servicename = "metadata"

type serviceConfig struct {
	APIConfig apiConfig `yaml:"api"`
}
type apiConfig struct {
	Port int `yaml:"port"`
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "api Handler port")
	flag.Parse()
	log.Println("Starting the movie metadata service")

	f, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	var cfg serviceConfig

	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		panic(err)
	}

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	instanceID := discovery.GenerateInstanceID(servicename)

	if err := registry.Register(ctx, instanceID, servicename, fmt.Sprintf("locahhost:%v", cfg.APIConfig.Port)); err != nil {
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
	defer registry.DeRegister(ctx, instanceID, servicename)
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.APIConfig.Port))
	if err != nil {
		log.Fatalf("faile to listen : %v", err.Error())
	}
	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, h)
	log.Println("Started the server on port: ", cfg.APIConfig.Port)
	srv.Serve(lis)
}
