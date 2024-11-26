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

type (
	GetAllPostResponse struct {
		Posts []Post `json:"posts"`
		Pagination Pagination `json:"pagination"`
	}

	Post struct {
		ID        int64       	`json:"id"`
		Username 	string				`json:"username"`
		Title     string    		`json:"title"`
		Content   string    		`json:"content"`
		Hashtags  []string    	`json:"hashtags"`
		IsLike		bool 					`json:"is_like"`
	}

	Pagination struct {
		Limit 				int `json:"limit"`
		Offset 				int `json:"offset"`
		Page					int `json:"page"`
		TotalPages		int `json:"total_pages"`
		TotalItems		int `json:"total_items"`
	}

	GetPostResponse struct {
		Post		Post `json:"post"`
		LikeCount 		int `json:"like_count"`
		Comments			[]Comment `json:"comments"`
	}

)