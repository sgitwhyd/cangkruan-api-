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
	commentRepo repository.CommentRepository
	userActRepo repository.UserActRepository
}

type PostService interface {
	Save(ctx context.Context, req model.CreatePostRequest, userID int64) error
	FindByID(ctx context.Context, userID, postID int64) (*model.GetPostResponse, error)
	FindAll(ctx context.Context, pageSize, pageIndex int) (model.GetAllPostResponse, error)
}

func NewPostService(repository repository.PostRepository, commentRepo repository.CommentRepository,
	userActRepo repository.UserActRepository) *service {
	return &service{repository: repository, commentRepo: commentRepo, userActRepo: userActRepo}
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

func (s *service) FindByID(ctx context.Context, userID, postID int64) (*model.GetPostResponse, error) {
	post, err := s.repository.FindByID(ctx,userID, postID)
	if err != nil {
		log.Error().Err(err).Msgf("service: post with id %d not found", postID)
		return nil, err
	}

	comments, err := s.commentRepo.GetCommentByPostID(ctx, postID)
	if err != nil {
		log.Error().Err(err).Msgf("service: comment with post_id %d not found", postID)
		return nil, err
	}

	like, err := s.userActRepo.CountLikeByID(ctx, postID)
	if err != nil {
		log.Error().Err(err).Msgf("service: like with post_id %d not found", postID)
		return nil, err
	}

	response := &model.GetPostResponse{
		Post: *post,
		LikeCount: int(like),
		Comments: comments,
	}

	return response, nil

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