package testutil

import (
	"bikraj.movie_microservice.net/gen"
	"bikraj.movie_microservice.net/metadata/internal/controller/metadata"
	grpchandler "bikraj.movie_microservice.net/metadata/internal/handler/grpc"
	"bikraj.movie_microservice.net/metadata/internal/repository/memory"
)

func NewTestMetadataGRPCServer() gen.MetadataServiceServer {
	r := memory.New()
	ctrl := metadata.New(r)
	return grpchandler.New(ctrl)
}
