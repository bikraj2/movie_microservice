package testutil

import (
	"bikraj.movie_microservice.net/gen"
	"bikraj.movie_microservice.net/movie/internal/controller"
	metadatagateway "bikraj.movie_microservice.net/movie/internal/gateway/metadata/grpc"
	ratinggateway "bikraj.movie_microservice.net/movie/internal/gateway/rating/grpc"
	grpchandler "bikraj.movie_microservice.net/movie/internal/handler/grpc"
	"bikraj.movie_microservice.net/pkg/discovery"
)

// NewTestMovieGRPCServer creates a new movie gRPC server to be used in tests.
func NewTestMovieGRPCServer(registry discovery.Registry) gen.MovieServiceServer {
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	return grpchandler.New(ctrl)
}
