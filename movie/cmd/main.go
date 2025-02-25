package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"bikraj.movie_microservice.net/gen"
	movie "bikraj.movie_microservice.net/movie/internal/controller"
	metadatagateway "bikraj.movie_microservice.net/movie/internal/gateway/metadata/grpc"
	ratinggateway "bikraj.movie_microservice.net/movie/internal/gateway/rating/grpc"
	grpchandler "bikraj.movie_microservice.net/movie/internal/handler/grpc"
	"bikraj.movie_microservice.net/pkg/discovery"
	consul "bikraj.movie_microservice.net/pkg/discovery/consul"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

const (
	serviceName = "movie"
)

type limiter struct {
	l *rate.Limiter
}

func newLimiter(limit int, burst int) *limiter {
	return &limiter{l: rate.NewLimiter(rate.Limit(limit), burst)}
}
func (l *limiter) Limit() bool {
	return l.l.Allow()
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Server listening on port %d", port)

	f, err := os.Open("../configs/base.yaml")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	var cfg serviceConfig

	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	registry, err := consul.NewRegistry("localhost:8500")

	if err != nil {
		panic(err)
	}

	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("locahost:%d", cfg.APIConfig.Port)); err != nil {
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

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", cfg.APIConfig.Port))
	if err != nil {
		log.Fatalf("err while listening: %v", err.Error())
	}

	const limit = 100
	const burst = 50
	l := newLimiter(limit, burst)
	srv := grpc.NewServer(grpc.UnaryInterceptor(ratelimit.UnaryServerInterceptor(l)))
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := <-sigChan
		cancel()
		log.Printf("Recieved singal :%v, attempting graceful shutdown", s)
		srv.GracefulStop()
		log.Println("Gracefully stoped the gRPC Server ")
	}()
	gen.RegisterMovieServiceServer(srv, h)
	srv.Serve(lis)
	wg.Wait()
}
