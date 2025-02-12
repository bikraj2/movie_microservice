package mysql

import (
	"context"
	"database/sql"

	model "bikraj.movie_microservice.net/metadata/pkg"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var movie model.Metadata

	query := `
  SELECT title,description, director
  FROM movie
  WHERE id = $1
  `
	err := r.db.QueryRowContext(ctx, query, id).Scan(&movie.Title, &movie.Description, &movie.Director)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *Repository) Put(ctx context.Context, id string, metata *model.Metadata) error {
	query := `
  INSERT INTO movie
  (id,title,description,director)
  VALUES ($1,$2,$3,$4)
  `
	_, err := r.db.ExecContext(ctx, query, id, metata.Title, metata.Description, metata.Director)
	return err
}
