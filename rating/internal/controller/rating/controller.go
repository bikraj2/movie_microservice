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
type ratingIngester interface {
	Ingest(ctx context.Context) (chan model.RatingEvent, error)
}

type Controller struct {
	repo ratingRepository

	ingester ratingIngester
}

func New(repo ratingRepository, ingester ratingIngester) *Controller {
	return &Controller{repo: repo, ingester: ingester}
}

// GetAgrregatedRaring returns the agrregated rating for
// a record or ErrRecordNotFound if not rating at all.

func (c *Controller) GetAgrregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
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

// StartIngestino starts the ingestion of rating events.
func (s *Controller) StartIngestion(ctx context.Context) error {
	ch, err := s.ingester.Ingest(ctx)

	if err != nil {
		return err
	}
	for e := range ch {
		if err := s.PutRating(ctx, e.RecordID, e.RecordType, &model.Rating{UserID: e.UserID, RatingValue: e.Value}); err != nil {
			return err
		}
	}
	return nil
}
