package model

import "time"

type (
	CommentModel struct {
		ID        int64     `db:"id"`
		PostID    int64     `db:"post_id"`
		UserID    int64     `db:"user_id"`
		Content   string    `db:"content"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
		CreatedBy string    `db:"created_by"`
		UpdatedBy string    `db:"updated_by"`
	}

	CreateCommentRequest struct {
		Content   string    `json:"content" binding:"required"`
	}
)