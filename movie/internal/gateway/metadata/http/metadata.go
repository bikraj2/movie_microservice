package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	model "bikraj.movie_microservice.net/metadata/pkg"
	"bikraj.movie_microservice.net/movie/internal/gateway"
	"bikraj.movie_microservice.net/pkg/discovery"
	"golang.org/x/exp/rand"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/metadata"
	log.Printf("Calling metadata service,Request GET:%v", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, gateway.ErrNotFound
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx respons: %v", resp)
	}
	var v *model.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
