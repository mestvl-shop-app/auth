package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mestvl-shop-app/auth/internal/app"
	"github.com/mestvl-shop-app/auth/internal/config"
	"github.com/mestvl-shop-app/auth/internal/db"
	"github.com/mestvl-shop-app/auth/internal/repository"
	"github.com/mestvl-shop-app/auth/internal/service"
	log "github.com/mestvl-shop-app/auth/pkg/logger"
)

func main() {
	// Init config
	cfg := config.MustLoad()

	// Init logger
	logger := log.SetupLogger(cfg.Env)
	logger.Info("start shop auth service",
		"env", cfg.Env,
	)
	logger.Debug("debug messages are enabled")

	// Init database
	dbPostgres, err := db.New(cfg.Database)
	if err != nil {
		logger.Error("postgres connect problem", "error", err)
		os.Exit(1)
	}
	defer func() {
		err = dbPostgres.Close()
		if err != nil {
			logger.Error("error when closing", "error", err)
		}
	}()
	logger.Info("postgres connection done")

	// Init repos, services, gRPC app
	repos := repository.NewRepositories(dbPostgres)
	services := service.NewServices(service.Deps{
		Logger: logger,
		Config: cfg,
		Repos:  repos,
	})

	app := app.New(
		logger,
		cfg,
		services,
	)

	go app.GRPCSrv.MustRun()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	logger.Info("stopping application",
		"signal",
		sign,
	)

	app.GRPCSrv.Stop()

	logger.Info("app stopped")
}
