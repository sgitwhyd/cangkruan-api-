package model

import "time"

type (
	SignUpRequest struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	SignInRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	SignInResponse struct {
		AccessToken string `json:"access_token"`
	}

	UserModel struct {
		ID        int64 `db:"id"`
		Username 	string `db:"username"`
		Email     string `db:"email"`
		Password  string `db:"password"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
		CreatedBy string `db:"created_by"`
		UpdatedBy string `db:"updated_by"`
	}
)