package metadata

import (
	"context"
	"errors"

	model "bikraj.movie_microservice.net/metadata/pkg"
)

var (
	ErrNotFound = errors.New("not found")
)

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
}

type Controller struct {
	repo metadataRepository
}

func New(repo metadataRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {

	res, err := c.repo.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return res, nil
}
