package service

import "errors"

var (
	ErrClientAlreadyExists = errors.New("client already exists")
	ErrClientNotFound      = errors.New("client not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrTokenExpired        = errors.New("token expired")
)
