package usecase

import (
	"context"
	"log/slog"
	"repo-stat/collector/internal/domain"
)

type CollectorProvider struct {
	log              *slog.Logger
	githubClient     GitHubClient
	subscriberClient SubscriberClient
}

func NewCollectorProvider(log *slog.Logger, githubClient GitHubClient, subscriberClient SubscriberClient) *CollectorProvider {
	return &CollectorProvider{
		log:              log,
		githubClient:     githubClient,
		subscriberClient: subscriberClient,
	}
}

func (u *CollectorProvider) GetRepo(ctx context.Context, owner, name string) (domain.Repo, error) {
	return u.githubClient.GetRepo(ctx, owner, name)
}

func (u *CollectorProvider) GetSubscriptionsInfo(ctx context.Context) ([]domain.Repo, error) {
	subscriptions, err := u.subscriberClient.GetSubscriptionsList(ctx)
	if err != nil {
		u.log.Error("Failed to get subscriptions", "error", err)
		return nil, err
	}

	result := make([]domain.Repo, 0, len(subscriptions))
	for _, sub := range subscriptions {
		repo, err := u.githubClient.GetRepo(ctx, sub.Owner, sub.Repo)
		if err != nil {
			u.log.Error("Failed to get repo", "error", err)
			continue
		}

		result = append(result, repo)
	}

	return result, nil
}
