package grpc

import (
	"context"

	"bikraj.movie_microservice.net/gen"
	grpcutil "bikraj.movie_microservice.net/internal/grpcutils"
	"bikraj.movie_microservice.net/pkg/discovery"
	ratingModel "bikraj.movie_microservice.net/rating/pkg/model"
)

type Gateway struct {
	registry discovery.Registry
}

func New(r discovery.Registry) *Gateway {
	return &Gateway{r}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordId ratingModel.RecordID, recordType ratingModel.RecordType) (float64, error) {
	conn, err := grpcutil.ServieConnect(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	client := gen.NewRatingServiceClient(conn)
	resp, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{RecordId: string(recordId), RecordType: string(recordType)})
	if err != nil {
		return 0, err
	}
	return resp.RatingValue, nil
}
func (g *Gateway) PutRating(ctx context.Context, recordId ratingModel.RecordID, recordType ratingModel.RecordType, rating *ratingModel.Rating) error {
	return nil
}
