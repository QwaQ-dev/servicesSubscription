package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/QwaQ-dev/servicesSubscription/internal/config"
	"github.com/QwaQ-dev/servicesSubscription/internal/handlers"
	postgres "github.com/QwaQ-dev/servicesSubscription/internal/repository"
	"github.com/QwaQ-dev/servicesSubscription/internal/routes"
	"github.com/QwaQ-dev/servicesSubscription/internal/services"
	"github.com/QwaQ-dev/servicesSubscription/pkg/sl"
	"github.com/gofiber/fiber/v2"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024,
	})

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("Starting subscriptions backend", slog.String("env", cfg.Env))
	db, err := postgres.InitDatabase(cfg.Database, log)
	if err != nil {
		log.Error("Error with connecting to database", sl.Err(err))
		os.Exit(1)
	}

	subscriptionRepo := postgres.NewSubsriptionRepo(db, log)
	subscriptionService := services.NewSubsriptionService(subscriptionRepo, log)
	subscriptionHandler := handlers.NewSubsriptionHandler(subscriptionService, log)

	routes.InitRoutes(app, log, subscriptionHandler)

	log.Info("starting server", slog.String("port", cfg.Server.Port))

	go func() {
		if err := app.Listen(cfg.Server.Port); err != nil {
			log.Error("Fiber server failed to start", sl.Err(err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	db.Close()

	log.Info("Shutting down application...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error("Error with shutting down Fiber server", sl.Err(err))
	} else {
		log.Info("Fiber server gracefully stopped.")
	}

	log.Info("Application exited.")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
