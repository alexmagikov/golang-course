package usecase

import (
	"context"
	"repo-stat/collector/internal/domain"
)

type GitHubClient interface {
	GetRepo(ctx context.Context, owner, name string) (domain.Repo, error)
}

type SubscriberClient interface {
	GetSubscriptionsList(ctx context.Context) ([]domain.Subscription, error)
}
