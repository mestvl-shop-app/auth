package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/mestvl-shop-app/auth/internal/service"
	authv1 "github.com/mestvl-shop-app/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
	services *service.Services
	logger   *slog.Logger
}

func Register(
	gRPC *grpc.Server,
	services *service.Services,
	logger *slog.Logger,
) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{
		services: services,
		logger:   logger,
	})
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.services.Auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) || errors.Is(err, service.ErrClientNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		s.logger.Error("login failed",
			"error", err,
			"email", req.GetEmail(),
		)
		return nil, status.Error(codes.Internal, "login failed")
	}

	return &authv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	clientId, err := s.services.Auth.Register(ctx, &service.RegisterDTO{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})

	if err != nil {
		if errors.Is(err, service.ErrClientAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		s.logger.Error("register new client failed",
			"error", err,
		)
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &authv1.RegisterResponse{
		UserId: clientId.String(),
	}, nil
}

func (s *serverAPI) Validate(ctx context.Context, req *authv1.ValidateRequest) (*authv1.ValidateResponse, error) {
	if err := s.services.Auth.ValidateToken(ctx, req.GetToken()); err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			return &authv1.ValidateResponse{
				Status: authv1.ValidateStatus_FORBIDDEN,
			}, nil
		}
		s.logger.Error("validate token failed",
			"error", err,
		)
		return nil, status.Error(codes.Internal, "validate token failed")
	}

	return &authv1.ValidateResponse{
		Status: authv1.ValidateStatus_OK,
	}, nil
}
