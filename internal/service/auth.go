package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mestvl-shop-app/auth/internal/domain"
	jwt_manager "github.com/mestvl-shop-app/auth/internal/pkg/jwt"
	"github.com/mestvl-shop-app/auth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	clientRepository repository.ClientInterface
	appRepository    repository.AppInterface
	logger           *slog.Logger
}

func newAuthService(
	clientRepository repository.ClientInterface,
	appRepository repository.AppInterface,
	logger *slog.Logger,
) *authService {
	return &authService{
		clientRepository: clientRepository,
		appRepository:    appRepository,
		logger:           logger,
	}
}

type RegisterDTO struct {
	Email    string
	Password string
}

func (s *authService) Register(ctx context.Context, dto *RegisterDTO) (*uuid.UUID, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt.GenerateFromPassword failed: %w", err)
	}

	clientId, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("uuid.NewV7 failed: %w", err)
	}

	if err := s.clientRepository.Create(ctx, &domain.Client{
		ID:       clientId,
		Email:    dto.Email,
		Password: passHash,
	}); err != nil {
		if errors.Is(err, domain.ErrDuplicateEntry) {
			return nil, ErrClientAlreadyExists
		}
		return nil, fmt.Errorf("create user failed: %w", err)
	}

	return &clientId, nil
}

func (s *authService) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	client, err := s.clientRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", ErrClientNotFound
		}
		return "", fmt.Errorf("get client by email failed: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword(client.Password, []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	app, err := s.appRepository.GetByID(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("get app by id failed: %w", err)
	}

	token, err := jwt_manager.NewToken(client, app)
	if err != nil {
		return "", fmt.Errorf("generate jwt token failed: %w", err)
	}

	return token, nil
}

func (s *authService) ValidateToken(ctx context.Context, accessToken string) error {
	appID, err := jwt_manager.GetAppID(accessToken)
	if err != nil {
		return fmt.Errorf("get app id from jwt failed: %w", err)
	}

	app, err := s.appRepository.GetByID(ctx, appID)
	if err != nil {
		return fmt.Errorf("get app by id failed: %w", err)
	}

	_, err = jwt_manager.Parse(accessToken, app.JwtSigningKey)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return ErrTokenExpired
		}
		return fmt.Errorf("parse token failed: %w", err)
	}

	return nil
}
