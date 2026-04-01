package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"repo-stat/collector/config"
	"repo-stat/collector/internal/adapter/github"
	"repo-stat/collector/internal/controller/grpcController"
	"repo-stat/collector/internal/usecase"
	"repo-stat/platform/logger"
	"repo-stat/proto/collector"

	"google.golang.org/grpc"
)

func run(ctx context.Context) error {
	// config
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	// logger
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting collector service...")
	log.Debug("debug messages are enabled")

	githubClient := github.NewClient(cfg.GitHub.Token)

	repoUseCase := usecase.NewRepoProvider(githubClient)

	repoServer := grpcController.NewRepoServer(repoUseCase)

	lis, err := net.Listen("tcp", cfg.GRPC.Address)
	if err != nil {
		log.Error("failed to listen", "address", cfg.GRPC.Address)
		return err
	}

	grpcServer := grpc.NewServer()
	collector.RegisterCollectorServiceServer(grpcServer, repoServer)
	
	log.Info("starting collector...")

	go func() {
		<-ctx.Done()
		log.Info("shutting down collector...")
		grpcServer.GracefulStop()
	}()

	return grpcServer.Serve(lis)
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching collector error: %s\n", err)
		}
		cancel()
		os.Exit(1)
	}
	cancel()
}
