package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	model "bikraj.movie_microservice.net/metadata/pkg"
	"bikraj.movie_microservice.net/movie/internal/gateway"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr: addr}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	req, err := http.NewRequest(http.MethodGet, g.addr+"/metadata", nil)
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
