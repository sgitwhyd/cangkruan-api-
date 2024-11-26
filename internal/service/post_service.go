package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/cangkruan-api/internal/model"
	"github.com/sgitwhyd/cangkruan-api/internal/repository"
)

type service struct {
	repository repository.PostRepository
}

type PostService interface {
	Save(ctx context.Context, req model.CreatePostRequest, userID int64) error
	FindByID(ctx context.Context, postID int64) (*model.PostModel, error)
	FindAll(ctx context.Context, pageSize, pageIndex int) (model.GetAllPostResponse, error)
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

func (s *service) FindByID(ctx context.Context, postID int64) (*model.PostModel, error) {
	post, err := s.repository.FindByID(ctx, postID)
	if err != nil {
		log.Error().Err(err).Msgf("post with id %d not found", postID)
		return nil, err
	}

	return post, nil

}

func (s *service) FindAll(ctx context.Context, pageSize, pageIndex int) (model.GetAllPostResponse, error) {
	limit := pageSize
	offset := pageSize * (pageIndex - 1)
	posts, err := s.repository.GetAll(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("service: error get all posts")
		return posts, err
	}

	return posts, err
}