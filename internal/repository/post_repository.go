package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/cangkruan-api/internal/model"
)

type postRepository struct {
	db *sql.DB
}

type PostRepository interface {
	Create(ctx context.Context, req model.PostModel) error
	FindByID(ctx context.Context, postID int64) (*model.PostModel, error)
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
		log.Error().Err(err).Msgf("repository error on create post title=%s, user_id=%d", req.Title, req.UserID)
		return err
	}

return nil

}

func (r *postRepository) FindByID(ctx context.Context, postID int64) (*model.PostModel, error) {
	query := `SELECT id, user_id, title, content, hashtags, created_at, updated_at, created_by, updated_by FROM posts WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, postID)

	var post model.PostModel
	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Hashtags, &post.CreatedAt, &post.UpdatedAt, &post.CreatedBy, &post.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post with ID %d not found", postID)
		}
		return nil, err
	}

	return &post, err
}