package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"bikraj.movie_microservice.net/gen"
	"bikraj.movie_microservice.net/pkg/discovery"
	"bikraj.movie_microservice.net/pkg/discovery/consul"
	"bikraj.movie_microservice.net/rating/internal/controller/rating"
	grpcHandler "bikraj.movie_microservice.net/rating/internal/handler/grpc"
	repo "bikraj.movie_microservice.net/rating/internal/repository/mysql"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

const (
	servicename = "rating"
)

type serviceConfig struct {
	APIConfig apiConfig `yaml:"api"`
}
type apiConfig struct {
	Port int `yaml:"port"`
}

func main() {

	var port int
	flag.IntVar(&port, "port", 8082, "api Handler port")
	flag.Parse()
	log.Println("Starting the movie metadata service")
	f, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	var cfg serviceConfig
	fmt.Println()
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
				log.Println("Failed to report healthy state: ", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.DeRegister(ctx, instanceID, servicename)
	db, err := sql.Open("mysql", "root:pa55word@/movie")
	if err != nil {
		panic(err)
	}
	repo := repo.New(db)
	ctrl := rating.New(repo, nil)
	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.APIConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err.Error())
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	reflection.Register(srv)
	log.Println("Started the server on port: ", cfg.APIConfig.Port)
	srv.Serve(lis)

}
