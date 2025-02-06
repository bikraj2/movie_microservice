package rating

import (
	"context"
	"errors"

	"bikraj.movie_microservice.net/rating/pkg/model"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type ratingRepository interface {
	Get(ctx context.Context, recordId model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordId model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

type Controller struct {
	repo ratingRepository
}

func New(repo ratingRepository) *Controller {
	return &Controller{repo: repo}
}

// GetAgrregatedRaring returns the agrregated rating for
// a record or ErrRecordNotFound if not rating at all.

func (c *Controller) GetAgrregatedRaring(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordID, recordType)
	if err != nil {
		switch {
		case errors.Is(err, ErrRecordNotFound):
			return 0, ErrRecordNotFound
		default:
			return 0, err
		}
	}
	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.RatingValue)
	}
	return sum / float64(len(ratings)), nil
}

func (c *Controller) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return c.repo.Put(ctx, recordID, recordType, rating)
}
