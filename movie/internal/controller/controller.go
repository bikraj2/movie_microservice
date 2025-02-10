package movie

import (
	"context"
	"errors"

	metadataModel "bikraj.movie_microservice.net/metadata/pkg"
	"bikraj.movie_microservice.net/movie/internal/gateway"
	"bikraj.movie_microservice.net/movie/pkg/model"
	ratingModel "bikraj.movie_microservice.net/rating/pkg/model"
)

var ErrNotFound = errors.New("not found")

type ratingGateway interface {
	PutRating(ctx context.Context, recordId ratingModel.RecordID, recordType ratingModel.RecordType, rating *ratingModel.Rating) error
	GetAggregatedRating(ctx context.Context, recordId ratingModel.RecordID, recordType ratingModel.RecordType) (float64, error)
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadataModel.Metadata, error)
}

type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

func New(rg ratingGateway, mg metadataGateway) *Controller {
	return &Controller{ratingGateway: rg, metadataGateway: mg}
}

func (ctrl *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {

	metadata, err := ctrl.metadataGateway.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	details := &model.MovieDetails{Metadata: *metadata}

	rating, err := ctrl.ratingGateway.GetAggregatedRating(ctx, ratingModel.RecordID(id), ratingModel.RecordTypeMovie)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
	} else if err != nil {
		details.Rating = &rating
	}
	return details, nil
}
