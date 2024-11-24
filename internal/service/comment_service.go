package service

import (
	"context"
	"strconv"
	"time"

	"github.com/sgitwhyd/cangkruan-api/internal/model"
	"github.com/sgitwhyd/cangkruan-api/internal/repository"
)

type commentService struct {
	repository repository.CommentRepository
}

type CommentService interface {
	Save(postID int64, ctx context.Context, req model.CreateCommentRequest, userID int64) error
}

func NewCommentService(repository repository.CommentRepository) *commentService { 
	return &commentService{repository: repository}
 }

func (s *commentService) 	Save(postID int64, ctx context.Context, req model.CreateCommentRequest, userID int64) error {

	now := time.Now()
	request := model.CommentModel{
		PostID: postID,
		UserID: userID,
		Content: req.Content,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: strconv.FormatInt(userID, 10),
		UpdatedBy: strconv.FormatInt(userID, 10),
	}

	err := s.repository.Create(ctx, request)
	if err != nil {
		return err
	}

	return err
}