package grpc

import (
	"context"

	"bikraj.movie_microservice.net/gen"
	grpcutil "bikraj.movie_microservice.net/internal/grpcutils"
	model "bikraj.movie_microservice.net/metadata/pkg"
	"bikraj.movie_microservice.net/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func New(r discovery.Registry) *Gateway {
	return &Gateway{r}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServieConnect(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataReqeust{MovieId: id})
	if err != nil {
		return nil, err
	}

	return model.MetadataFromProto(resp.Metadata), nil
}
