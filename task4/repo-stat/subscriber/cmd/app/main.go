package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	subscriberpb "repo-stat/proto/subscriber"
	"repo-stat/subscriber/config"
	"repo-stat/subscriber/internal/adapter/github"
	grpccontroller "repo-stat/subscriber/internal/controller/grpc"
	"repo-stat/subscriber/internal/repository/postgres"
	"repo-stat/subscriber/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting subscriber server...")
	log.Debug("debug messages are enabled")

	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	))
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		return err
	}
	defer pool.Close()

	queries := postgres.New(pool)

	githubClient := github.NewClient(cfg.Github.Token)

	pingUseCase := usecase.NewPing()
	subscriptionProviderUsecase := usecase.NewSubscriptionProvider(log, githubClient, queries)
	server := grpccontroller.NewServer(log, pingUseCase, subscriptionProviderUsecase)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	subscriberpb.RegisterSubscriberServiceServer(srv.GRPC(), server)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching server error: %s\n", err)
		}
		cancel()
		os.Exit(1)
	}
}
