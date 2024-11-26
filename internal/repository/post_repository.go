package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/cangkruan-api/internal/model"
)

type postRepository struct {
	db *sql.DB
}

type PostRepository interface {
	GetAll(ctx context.Context, limit int, offset int)(model.GetAllPostResponse, error) 
	Create(ctx context.Context, req model.PostModel) error
	FindByID(ctx context.Context, postID int64) (*model.PostModel, error)
}

func NewPostRepository(db *sql.DB) *postRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) GetAll(ctx context.Context, limit int, offset int)(model.GetAllPostResponse, error) {
	var posts model.GetAllPostResponse
	query := `SELECT p.id, p.user_id, p.title, p.content, p.hashtags, p.created_at, p.updated_at, p.created_by, p.updated_by, u.username FROM posts p JOIN users u ON p.user_id = u.id ORDER BY p.updated_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return posts, err
	}

	query = `SELECT COUNT(*) as total_rows FROM posts`
	var totalItems int64
	err = r.db.QueryRowContext(ctx, query).Scan(&totalItems)
	if err != nil {
		return posts, err
	}

	defer rows.Close()

	result := make([]model.Post, 0)
	for rows.Next() {
		var (
			post model.PostModel
			username string
		)
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Hashtags, &post.CreatedAt, &post.UpdatedAt, &post.CreatedBy, &post.UpdatedBy, &username)
		if err != nil {
			if err != sql.ErrNoRows {
				return posts, err
			}

			return posts, err
		}

		result = append(result, model.Post{
			ID: post.ID,
			Title: post.Title,
			UserID: post.UserID,
			Username: username,
			Content: post.Content,
			Hashtags: strings.Split(post.Hashtags, ","),
		})
	}

	posts.Posts = result
	posts.Pagination = model.Pagination{
		Limit: limit,
		Offset: offset,
		Page: (offset / limit) + 1,
		TotalItems: int(totalItems),
		TotalPages: (int(totalItems) + limit - 1) / limit,
	}
	
	return posts, nil
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