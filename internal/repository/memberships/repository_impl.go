package memberships

import (
	"context"
	"database/sql"

	memberships "github.com/sgitwhyd/cangkruan-api/internal/model"
)

type Repository interface {
	GetUser(ctx context.Context, email string, username string) (*memberships.UserModel, error)
	CreateUser(ctx context.Context, model memberships.UserModel) error
}

func (r *repository) GetUser(ctx context.Context, email string, username string) (*memberships.UserModel, error) {
	query := `SELECT id, email, password, username, created_at, updated_at, created_by, updated_by   FROM users WHERE email = ? OR username = ?` 
	rows := r.db.QueryRowContext(ctx, query, email, username)

	var user memberships.UserModel
	err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.CreatedAt, &user.UpdatedAt, &user.CreatedBy, &user.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) CreateUser(ctx context.Context, model memberships.UserModel) error {
	query := `INSERT INTO users (email, password, username, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, model.Email, model.Password, model.Username, model.CreatedAt, model.UpdatedAt, model.CreatedBy, model.UpdatedBy)

	if err != nil {
		return err
	}

	return nil
	
}