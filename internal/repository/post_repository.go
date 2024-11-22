package repository

import (
	"context"
	"database/sql"

	"github.com/sgitwhyd/cangkruan-api/internal/model"
)

type postRepository struct {
	db *sql.DB
}

type PostRepository interface {
	Create(ctx context.Context, req model.PostModel) error
}

func NewPostRepository(db *sql.DB) *postRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) Create(ctx context.Context, req model.PostModel) error {
	

	query := `INSERT INTO posts (title, content, hashtags, user_id, created_by, updated_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	_, err := r.db.ExecContext(ctx, query, req.Title, req.Content, req.Hashtags, req.UserID, req.CreatedBy, req.UpdatedBy, req.CreatedAt, req.UpdatedAt)

	if err != nil {
		return err
	}

return nil

}