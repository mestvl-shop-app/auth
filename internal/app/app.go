package app

import (
	"log/slog"

	grpcapp "github.com/mestvl-shop-app/auth/internal/app/grpc"
	"github.com/mestvl-shop-app/auth/internal/config"
	"github.com/mestvl-shop-app/auth/internal/service"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	services *service.Services,
) *App {
	grpcApp := grpcapp.New(log, cfg, services)

	return &App{
		GRPCSrv: grpcApp,
	}
}
