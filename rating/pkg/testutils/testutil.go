package testutil

import (
	"bikraj.movie_microservice.net/gen"
	"bikraj.movie_microservice.net/rating/internal/controller/rating"
	grpchandler "bikraj.movie_microservice.net/rating/internal/handler/grpc"
	"bikraj.movie_microservice.net/rating/internal/repository/memory"
)

// NewTestRatingGRPCServer creates a new rating gRPC server to be used in tests.
func NewTestRatingGRPCServer() gen.RatingServiceServer {
	r := memory.New()
	ctrl := rating.New(r, nil)
	return grpchandler.New(ctrl)
}
