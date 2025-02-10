package grpc

import (
	"bikraj.movie_microservice.net/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func New(r discovery.Registry) *Gateway {
	return &Gateway{r}
}


