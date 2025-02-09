package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"bikraj.movie_microservice.net/movie/internal/gateway"
	"bikraj.movie_microservice.net/pkg/discovery"
	"bikraj.movie_microservice.net/rating/pkg/model"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

func (g *Gateway) GetAggreatedRating(ctx context.Context, recordId model.RecordID, recordType model.RecordType) (float64, error) {
	addr, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return 0, err
	}
	if len(addr) == 0 {
		return 0, err
	}
	url := "http://" + addr[rand.Intn(len(addr))] + "/rating"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)
	value := req.URL.Query()
	value.Add("id", string(recordId))
	value.Add("type", string(recordType))
	req.URL.RawQuery = value.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("not 2xx-code: %v", err)
	}

	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, nil
	}
	return v, nil
}

func (g *Gateway) PutRating(ctx context.Context, recordId model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	addr, err := g.registry.ServiceAddresses(ctx, "rating")
	if err != nil {
		return err
	}
	if len(addr) == 0 {
		return err
	}
	url := "http://" + addr[rand.Intn(len(addr))] + "/rating"
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordId))
	values.Add("type", string(recordType))
	values.Add("userId", string(rating.UserID))

	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non 2xx-code: %v", resp)
	}
	return nil
}
