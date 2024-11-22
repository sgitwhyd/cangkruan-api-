package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/sgitwhyd/cangkruan-api/internal/model"
	"github.com/sgitwhyd/cangkruan-api/internal/repository"
)

type service struct {
	repository repository.PostRepository
}

type PostService interface {
	Save(ctx context.Context, req model.CreatePostRequest, userID int64) error
}

func NewPostService(repository repository.PostRepository) *service {
	return &service{repository: repository}
}

func (s *service) Save(ctx context.Context, req model.CreatePostRequest, userID int64) error {

		hashtags := strings.Join(req.Hashtags, ",")

	now := time.Now()
	request := &model.PostModel{
		UserID: userID,
		Title: req.Title,
		Content: req.Content,
		Hashtags: hashtags,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: strconv.FormatInt(userID, 10),
		UpdatedBy: strconv.FormatInt(userID, 10),
	}

	err := s.repository.Create(ctx, *request)
	if err != nil {
		return err
	}

	return nil
}