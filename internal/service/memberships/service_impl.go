package memberships

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	model "github.com/sgitwhyd/cangkruan-api/internal/model"
	"github.com/sgitwhyd/cangkruan-api/pkg/jwt"
	tokenUtils "github.com/sgitwhyd/cangkruan-api/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) SignUp(ctx context.Context, req model.SignUpRequest) error {
	user, err := s.repository.GetUser(ctx, req.Email, req.Username, 0)
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

func (s *service) SignIn(ctx context.Context, req model.SignInRequest) (string, string, error) {
	user, err := s.repository.GetUser(ctx, req.Email, "", 0)
	if err != nil {
		log.Error().Err(err).Msg("failed get user")
		return "", "", nil
	}

	if user == nil {
		return "", "", errors.New("email not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", "", errors.New("email or password is invalid")
	}

	token, err := jwt.CreateToken(user.ID, user.Username, s.cfg.Service.SecretJWT)
	if err != nil {
		return "", "", err
	}

	existingRefToken, err := s.repository.GetRefreshToken(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Msgf("service: error getting latest refresh token by user_id: %d", user.ID)
		return "", "", err
	}

	if existingRefToken != nil {
		return token, existingRefToken.RefreshToken, nil
	}

	refreshToken := tokenUtils.GenerateToken()
	if refreshToken == "" {
		return token, "", errors.New("failed generate refresh token")
	}

	now := time.Now()
	err = s.repository.CreateRefreshToken(ctx, model.RefreshTokenModel{
		UserID: user.ID,
		RefreshToken: refreshToken,
		ExpiredAt: now.Add(10 * 24 * time.Hour),
		CreatedAt: now,
		Updatedat: now,
	})

	if err != nil {
		log.Error().Err(err).Msgf("service: error insert refresh token user_id: %d", user.ID)
		return token, "", err
	}

	return token, refreshToken, nil
}

func (s *service) ValidateRefreshToken(ctx context.Context, userID int64, request model.RefreshTokenRequest) (string, error) {
	exRefreshToken, err := s.repository.GetRefreshToken(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msgf("service: failed get refresh token user_id: %d", userID)
		return "", err
	}

	if exRefreshToken == nil {
		return "", errors.New("refresh token has expired")
	}

	if request.Token != exRefreshToken.RefreshToken {
		return "", errors.New("refresh token is invalid")
	}

	user, err := s.repository.GetUser(ctx, "", "", userID)
	if err != nil {
		log.Error().Err(err).Msg("failed get user")
		return "", nil
	}

	accessToken, err := jwt.CreateToken(userID, user.Username, s.cfg.Service.SecretJWT)
	if err != nil {
		log.Error().Err(err).Msgf("service: failed generate jwt token user_id: %d", userID)
		return  "", err
	}

	return accessToken, nil

}