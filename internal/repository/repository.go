package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/mestvl-shop-app/auth/internal/domain"
)

type Repositories struct {
	Client ClientInterface
	App    AppInterface
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Client: newClientRepository(db),
		App:    newAppRepository(db),
	}
}

type ClientInterface interface {
	Create(ctx context.Context, client *domain.Client) error
	GetByEmail(ctx context.Context, email string) (*domain.Client, error)
}

type AppInterface interface {
	GetByID(ctx context.Context, id int) (*domain.App, error)
}
