package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mestvl-shop-app/auth/internal/domain"
)

type appRepository struct {
	db *sqlx.DB
}

func newAppRepository(db *sqlx.DB) *appRepository {
	return &appRepository{db: db}
}

func (r *appRepository) GetByID(ctx context.Context, id int) (*domain.App, error) {
	const query = `
	SELECT id, "name", jwt_signing_key, jwt_access_token_ttl_minutes, jwt_refresh_token_ttl_minutes, created_at, updated_at, deleted_at
	FROM app
	WHERE id = $1;
	`

	var app domain.App
	if err := r.db.GetContext(ctx, &app, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("select app failed: %w", err)
	}

	return &app, nil
}
