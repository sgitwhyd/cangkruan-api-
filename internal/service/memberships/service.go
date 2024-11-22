package memberships

import (
	"context"

	"github.com/sgitwhyd/cangkruan-api/internal/configs"
	model "github.com/sgitwhyd/cangkruan-api/internal/model"
	"github.com/sgitwhyd/cangkruan-api/internal/repository/memberships"
)

type service struct {
	cfg *configs.Config
	repository memberships.Repository
}

type MembershipService interface {
	SignUp(ctx context.Context, req model.SignUpRequest) error
	SignIn(ctx context.Context, req model.SignInRequest) (string, error)
}

func NewService(cfg *configs.Config,membershipRepo memberships.Repository) *service {
	return &service{cfg, membershipRepo}
}