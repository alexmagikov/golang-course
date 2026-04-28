package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type ProcessorClient interface {
	Ping(ctx context.Context) domain.PingStatus
	GetRepo(ctx context.Context, url string) (domain.Repo, error)
	GetSubscriptionInfo(ctx context.Context) ([]domain.Repo, error)
	Close() error
}

type SubscriberClient interface {
	Ping(ctx context.Context) domain.PingStatus
	CreateSubscription(ctx context.Context, url string) error
	DeleteSubscription(ctx context.Context, owner, repo string) error
	GetSubscriptionsList(ctx context.Context) ([]domain.Subscription, error)
	Close() error
}
