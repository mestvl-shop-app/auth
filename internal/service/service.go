package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mestvl-shop-app/auth/internal/config"
	"github.com/mestvl-shop-app/auth/internal/repository"
)

type Services struct {
	Auth AuthInterface
}

type Deps struct {
	Logger *slog.Logger
	Config *config.Config
	Repos  *repository.Repositories
}

func NewServices(deps Deps) *Services {
	return &Services{
		Auth: newAuthService(
			deps.Repos.Client,
			deps.Repos.App,
			deps.Logger,
		),
	}
}

type AuthInterface interface {
	Register(ctx context.Context, dto *RegisterDTO) (*uuid.UUID, error)
	Login(ctx context.Context, email string, password string, appID int) (string, error)
	ValidateToken(ctx context.Context, accessToken string) error
}
