package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/cangkruan-api/internal/model"
	"github.com/sgitwhyd/cangkruan-api/internal/repository"
)

type userActService struct {
	repository repository.UserActRepository
}

type UserActService interface{
	Create(ctx context.Context, req model.CreateUserActivityRequest, userID int64, postID int64) error
}

func NewUserActivityService(repository repository.UserActRepository) *userActService {
	return &userActService{
		repository: repository,
	}
}

func (s *userActService) 	Create(ctx context.Context, req model.CreateUserActivityRequest, userID int64, postID int64) error {	

	isLiked := 0
	if req.IsLiked {
		isLiked = 1
	}
	
	now := time.Now()
	model := model.UserActivityModel{
		UserID: userID,
		PostID: postID,
		IsLiked: isLiked,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: strconv.FormatInt(userID, 10),
		UpdatedBy: strconv.FormatInt(userID, 10),
	}

	userAct, err := s.repository.Find(ctx, model)
	
	if userAct == nil {

		if !req.IsLiked {
			return errors.New("you haven't liked this post yet")
		}

		// create
		err = s.repository.Create(ctx, model)
	} else {
		// update
		err = s.repository.Update(ctx, model)
	}

	if err != nil {
		log.Error().Err(err).Msgf("error like post on PostID: %d, UserID: %d", model.PostID, model.UserID)
		return err
	}


	return nil

}
