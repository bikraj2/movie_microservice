package mysql

import (
	"context"
	"database/sql"

	repository "bikraj.movie_microservice.net/rating/internal/repository"
	"bikraj.movie_microservice.net/rating/pkg/model"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(ctx context.Context, record_id model.RecordID, record_type model.RecordType) ([]model.Rating, error) {
	var res []model.Rating

	query := ` SELECT user_id,value
  FROM ratings
  WHERE  record_id = ? and record_type = ?
  `
	rows, err := r.db.QueryContext(ctx, query, record_id, record_type)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var rating model.Rating
		err := rows.Scan(&rating.UserID, &rating.RatingValue)
		if err != nil {
			return nil, err
		}

		res = append(res, rating)
	}
	if len(res) == 0 {
		return nil, repository.ErrRecordNotFound
	}
	return res, nil
}

func (r *Repository) Put(ctx context.Context, record_id model.RecordID, record_type model.RecordType, rating *model.Rating) error {
	query := `
  INSERT INTO ratings
  (record_id,record_type,user_id,value)
  VALUES (?,?,?,?)
  `
	_, err := r.db.ExecContext(ctx, query, record_id, record_type, rating.UserID, rating.RatingValue)
	return err
}
