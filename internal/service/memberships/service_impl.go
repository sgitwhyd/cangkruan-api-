package memberships

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	model "github.com/sgitwhyd/cangkruan-api/internal/model"
	"github.com/sgitwhyd/cangkruan-api/pkg/jwt"
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

func (s *service) SignIn(ctx context.Context, req model.SignInRequest) (string, error) {
	user, err := s.repository.GetUser(ctx, req.Email, "")
	if err != nil {
		log.Error().Err(err).Msg("failed get user")
		return "", nil
	}

	if user == nil {
		return "", errors.New("email not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("email or password is invalid")
	}

	token, err := jwt.CreateToken(user.ID, user.Username, s.cfg.Service.SecretJWT)
	if err != nil {
		return "", err
	}

	return token, nil


}