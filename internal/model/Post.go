package model

import "time"

type (
	PostModel struct {
		ID        int64       `db:"id"`
		UserID    int64       `db:"user_id"`
		Title     string    `db:"title"`
		Content   string    `db:"content"`
		Hashtags   string    `db:"hashtags"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
		CreatedBy string    `db:"created_by"`
		UpdatedBy string    `db:"updated_by"`
	}

	CreatePostRequest struct {
		Title     string    `json:"title" binding:"required"`
		Content   string    `json:"content" binding:"required"`
		Hashtags    []string    `json:"hashtags" binding:"required"`
	}
)