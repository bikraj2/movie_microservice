package grpc

import (
	"context"

	"bikraj.movie_microservice.net/gen"
	grpcutil "bikraj.movie_microservice.net/internal/grpcutils"
	model "bikraj.movie_microservice.net/metadata/pkg"
	"bikraj.movie_microservice.net/pkg/discovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	var resp *gen.GetMetadataResponse
	const maxRetries = 5
	for range 5 {
		resp, err = client.GetMetadata(ctx, &gen.GetMetadataReqeust{MovieId: id})
		if err != nil {
			if shouldRetry(err) {
				continue
			}
			return nil, err
		}
		return model.MetadataFromProto(resp.Metadata), nil
	}
	if err != nil {
		return nil, err
	}

	return model.MetadataFromProto(resp.Metadata), nil
}

func shouldRetry(err error) bool {

	e, ok := status.FromError(err)
	if !ok {
		return false
	}
	return e.Code() == codes.DeadlineExceeded || e.Code() == codes.ResourceExhausted || e.Code() == codes.Unavailable
}
