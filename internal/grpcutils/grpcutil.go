package grpcutil

import (
	"context"
	"math/rand"

	"bikraj.movie_microservice.net/pkg/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServieConnect(ctx context.Context, serviceName string, registry discovery.Registry) (*grpc.ClientConn, error) {
	addrs, err := registry.ServiceAddresses(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(addrs[rand.Intn(len(addrs))], grpc.WithTransportCredentials(insecure.NewCredentials()))
}
