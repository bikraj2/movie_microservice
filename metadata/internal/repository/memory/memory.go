package memory

import (
	"context"
	"sync"

	"bikraj.movie_microservice.net/metadata/internal/repository"
	model "bikraj.movie_microservice.net/metadata/pkg"
)

// In memory implemntation of our database.
// RWMutex to sync the any type of operations on the map data.

type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

// Creates a In Memory Repository.
func New() *Repository {
	return &Repository{data: map[string]*model.Metadata{}}
}

func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {

	r.Lock()
	defer r.Unlock()
	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

func (r *Repository) Put(_ context.Context, id string, metadata *model.Metadata) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}
