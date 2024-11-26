package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/cangkruan-api/internal/model"
)

type CommentRepository interface {
	Create(ctx context.Context, model model.CommentModel) error
	GetCommentByPostID(ctx context.Context, postID int64) ([]model.Comment, error)
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

func (r commentRepository) GetCommentByPostID(ctx context.Context, postID int64) ([]model.Comment, error) {

	comments := make([]model.Comment, 0)
	query := `SELECT c.id, c.content, u.username  FROM comments c JOIN users u ON c.user_id = u.id WHERE c.post_id = ?`

	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			return comments, fmt.Errorf("comment not found")
		}
		return comments, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			comment model.CommentModel
			username string
		)

		err := rows.Scan(&comment.ID, &comment.Content, &username)
		if err != nil {
			if err != sql.ErrNoRows {
				return comments, err
			}

			return comments, err
		}

		comments = append(comments, model.Comment{
			ID: comment.ID,
			Content: comment.Content,
			CreatedBy: username,
		})
	}

	return comments, nil
}
