package memberships

import (
	"context"
	"errors"
	"log"
	"time"

	model "github.com/sgitwhyd/cangkruan-api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) SignUp(ctx context.Context, req model.SignUpRequest) error {
	user, err := s.repository.GetUser(ctx, req.Email, req.Username)
	if err != nil {
		log.Printf("error create user %+v", err)
		return err
	}

	if user != nil {
		return errors.New("username or email already exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	model := model.UserModel{
		Username: req.Username,
		Email    : req.Email,
		Password : string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: req.Email,
		UpdatedBy: req.Email,
	}

	err = s.repository.CreateUser(ctx, model)
	if err != nil {
		return err
	}

	return nil

}