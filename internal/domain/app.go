package domain

import (
	"time"
)

type App struct {
	ID                        int        `db:"id"`
	Name                      string     `db:"name"`
	JwtSigningKey             string     `db:"jwt_signing_key"`
	JwtAccessTokenTtlMinutes  int        `db:"jwt_access_token_ttl_minutes"`
	JwtRefreshTokenTtlMinutes int        `db:"jwt_refresh_token_ttl_minutes"`
	CreatedAt                 time.Time  `db:"created_at"`
	UpdatedAt                 time.Time  `db:"updated_at"`
	DeletedAt                 *time.Time `db:"deleted_at"`
}
