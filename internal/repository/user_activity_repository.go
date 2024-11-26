package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/cangkruan-api/internal/model"
)

type userActRepository struct {
	db *sql.DB
}

type UserActRepository interface {
	Find(ctx context.Context, userAct model.UserActivityModel) (*model.UserActivityModel, error)
	Create(ctx context.Context, userAct model.UserActivityModel) error
	Update(ctx context.Context, userAct model.UserActivityModel) error
	CountLikeByID(ctx context.Context, postID int64) (int64, error)
}

func NewUserActivityRepository(db *sql.DB) *userActRepository {
	return &userActRepository{db: db}
}

func (r *userActRepository) Find(ctx context.Context, userAct model.UserActivityModel) (*model.UserActivityModel, error) {
	query := `SELECT id, post_id, user_id, is_liked, created_by, updated_by, created_at, updated_at FROM user_activities WHERE user_id = ? AND post_id = ?`
	row := r.db.QueryRowContext(ctx, query, userAct.UserID, userAct.PostID)
	
	var response model.UserActivityModel
	err := row.Scan(&userAct.ID, &userAct.PostID, &userAct.UserID, &userAct.IsLiked, &userAct.CreatedBy, &userAct.UpdatedBy, &userAct.CreatedAt, &userAct.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msgf("user activity with user_id: %d not found, post_id: %d", userAct.UserID, userAct.PostID)
			return nil, fmt.Errorf("user activity with user_id: %d not found", userAct.UserID)
		}
		return nil, err
	}

	return &response, nil
}
	
func (r *userActRepository) Create(ctx context.Context, userAct model.UserActivityModel) error {
	query := `INSERT INTO user_activities (post_id, user_id, is_liked, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?);`
	_, err := r.db.ExecContext(ctx, query, userAct.PostID, userAct.UserID, userAct.IsLiked, userAct.CreatedAt, userAct.UpdatedAt, userAct.CreatedBy, userAct.UpdatedBy)
	if err != nil {
			log.Error().Err(err).Msgf("repository error on create user activity, user_id: %d, post_id: %d", userAct.UserID, userAct.PostID)
		return err
	}

	return nil

}

func (r *userActRepository) Update(ctx context.Context, userAct model.UserActivityModel) error {
	query := `UPDATE user_activities SET is_liked = ?, created_by = ?, updated_by = ? WHERE post_id = ? AND user_id = ?`
	d, err := r.db.ExecContext(ctx, query, userAct.IsLiked, userAct.CreatedBy, userAct.UpdatedBy , userAct.PostID, userAct.UserID)
	if err != nil {
		log.Error().Err(err).Msgf("repository: failed update post with id:%d, user_id:%d", userAct.PostID, userAct.UserID)
		return err
	}

	rowsAffected, err := d.RowsAffected()
	if err != nil {
		log.Printf("repository: unable to fetch rows affected for user activity update (PostID: %d, UserID: %d). Error: %v", userAct.PostID, userAct.UserID, err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("repository: no rows updated for user activity (PostID: %d, UserID: %d)", userAct.PostID, userAct.UserID)
		return nil // Not an error, but indicates nothing was updated
	}

	return nil
}

func (r *userActRepository) CountLikeByID(ctx context.Context, postID int64) (int64, error) {
	query := `SELECT count(*) AS "like" FROM user_activities WHERE post_id = ? AND is_liked = 1`

	var like int64
	row := r.db.QueryRowContext(ctx, query, postID)
	
	err := row.Scan(&like)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, err
		}
		return 0, err
	}

	return like, nil
}