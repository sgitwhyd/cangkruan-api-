package repository

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/cangkruan-api/internal/model"
)

type CommentRepository interface {
	Create(ctx context.Context, model model.CommentModel) error
}

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *commentRepository {
	return &commentRepository{db: db}
}

func (r commentRepository) Create(ctx context.Context, model model.CommentModel)  error {
	query := `INSERT INTO comments (post_id, user_id, content, created_by, updated_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?);`
	_, err := r.db.ExecContext(ctx, query, model.PostID, model.UserID, model.Content, model.CreatedBy, model.UpdatedBy, model.CreatedAt, model.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Msgf("error create comment on post ID=%d, user_id=%d", model.PostID, model.UserID)
		return err
	}

	return nil
}
