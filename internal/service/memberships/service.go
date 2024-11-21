package memberships

import (
	"context"

	model "github.com/sgitwhyd/cangkruan-api/internal/model/memberships"
	"github.com/sgitwhyd/cangkruan-api/internal/repository/memberships"
)

type service struct {
	repository memberships.Repository
}

type MembershipService interface {
	SignUp(ctx context.Context, req model.SignUpRequest) error
}

func NewService(membershipRepo memberships.Repository) *service {
	return &service{membershipRepo}
}