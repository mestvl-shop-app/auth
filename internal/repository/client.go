package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mestvl-shop-app/auth/internal/db"
	"github.com/mestvl-shop-app/auth/internal/domain"
)

type clientRepository struct {
	db *sqlx.DB
}

func newClientRepository(db *sqlx.DB) *clientRepository {
	return &clientRepository{
		db: db,
	}
}

func (r *clientRepository) Create(ctx context.Context, client *domain.Client) error {
	const query = `
	INSERT INTO client
	(id, email, password, created_at, updated_at)
	VALUES(:id, :email, :password, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
	`

	_, err := r.db.NamedExecContext(ctx, query, client)
	if err != nil {
		if db.IsDuplicate(err) {
			return domain.ErrDuplicateEntry
		}

		return fmt.Errorf("insert client failed: %w", err)
	}

	return nil
}

func (r *clientRepository) GetByEmail(ctx context.Context, email string) (*domain.Client, error) {
	const query = `
	SELECT id, email, password, created_at, updated_at, deleted_at
	FROM client 
	WHERE email = $1;
	`

	var client domain.Client
	if err := r.db.GetContext(ctx, &client, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("select client failed: %w", err)
	}

	return &client, nil
}
