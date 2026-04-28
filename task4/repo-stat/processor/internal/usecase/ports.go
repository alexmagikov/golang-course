package usecase

import (
	"context"
	"repo-stat/processor/internal/domain"
)

type CollectorClient interface {
	GetRepo(ctx context.Context, owner, name string) (domain.Repo, error)
	GetSubscriptionInfo(ctx context.Context) ([]domain.Repo, error)
	Close() error
}
