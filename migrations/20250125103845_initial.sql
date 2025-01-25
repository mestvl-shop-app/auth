-- +goose Up
-- +goose StatementBegin
CREATE TABLE client (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE app (
    id SMALLINT PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    jwt_signing_key VARCHAR(255) NOT NULL,
    jwt_access_token_ttl_minutes INT NOT NULL,
    jwt_refresh_token_ttl_minutes INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

INSERT INTO app (id, name, jwt_signing_key, jwt_access_token_ttl_minutes, jwt_refresh_token_ttl_minutes) VALUES (1, 'shop', 'notasecret', 60, 120);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE client;
DROP TABLE app;
-- +goose StatementEnd
