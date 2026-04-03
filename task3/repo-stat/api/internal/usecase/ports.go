package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type ProcessorClient interface {
	Ping(ctx context.Context) domain.PingStatus
	GetRepo(ctx context.Context, url string) (domain.Repo, error)
	Close() error
}

type SubscriberClient interface {
	Ping(ctx context.Context) domain.PingStatus
	Close() error
}
