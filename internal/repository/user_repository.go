package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/sgitwhyd/cangkruan-api/internal/model"
)


type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}


type UserRepository interface {
	GetUser(ctx context.Context, email string, username string, userID int64) (*model.UserModel, error)
	CreateUser(ctx context.Context, model model.UserModel) error
	CreateRefreshToken(ctx context.Context, model model.RefreshTokenModel) error
	GetRefreshToken(ctx context.Context, userID int64) (*model.RefreshTokenModel, error)
}

func (r *userRepository) GetUser(ctx context.Context, email string, username string, userID int64) (*model.UserModel, error) {
	query := `SELECT id, email, password, username, created_at, updated_at, created_by, updated_by   FROM users WHERE email = ? OR username = ? OR id = ?` 
	rows := r.db.QueryRowContext(ctx, query, email, username, userID)

	var user model.UserModel
	err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.CreatedAt, &user.UpdatedAt, &user.CreatedBy, &user.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, model model.UserModel) error {
	query := `INSERT INTO users (email, password, username, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, model.Email, model.Password, model.Username, model.CreatedAt, model.UpdatedAt, model.CreatedBy, model.UpdatedBy)

	if err != nil {
		return err
	}

	return nil
	
}

func (r *userRepository) CreateRefreshToken(ctx context.Context, model model.RefreshTokenModel) error {
	query := `INSERT INTO refresh_token (user_id, refresh_token, expired_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, model.UserID, model.RefreshToken, model.ExpiredAt, model.CreatedAt, model.Updatedat)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetRefreshToken(ctx context.Context, userID int64) (*model.RefreshTokenModel, error) {
	query := `SELECT id, user_id, refresh_token, expired_at, created_at, updated_at FROM refresh_token WHERE user_id = ? AND expired_at >= ?`

	now := time.Now()
	var model model.RefreshTokenModel
	err := r.db.QueryRowContext(ctx, query, userID, now).Scan(&model.ID, &model.UserID, &model.RefreshToken, &model.ExpiredAt, &model.CreatedAt, &model.Updatedat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &model, nil
}