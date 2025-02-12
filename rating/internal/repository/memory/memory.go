package memory

import (
	"context"
	"errors"

	"bikraj.movie_microservice.net/rating/pkg/model"
)

type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

func New() *Repository {
	return &Repository{make(map[model.RecordType]map[model.RecordID][]model.Rating)}
}

func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {

	if _, ok := r.data[recordType]; !ok {
		return nil, ErrRecordNotFound
	}

	if ratings, ok := r.data[recordType][recordID]; !ok || len(ratings) == 0 {
		return nil, ErrRecordNotFound
	}
	return r.data[recordType][recordID], nil
}

func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {

	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = make(map[model.RecordID][]model.Rating)

	}
	r.data[recordType][recordID] = append(r.data[recordType][recordID], *rating)
	return nil
}
