package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type SubscriptionProvider struct {
	subscriberClient SubscriberClient
	processorClient  ProcessorClient
}

func NewSubscriptionProvider(subscriberClient SubscriberClient, processorClient ProcessorClient) *SubscriptionProvider {
	return &SubscriptionProvider{
		subscriberClient: subscriberClient,
		processorClient:  processorClient,
	}
}

func (p *SubscriptionProvider) CreateSubscription(ctx context.Context, url string) error {
	return p.subscriberClient.CreateSubscription(ctx, url)
}

func (p *SubscriptionProvider) DeleteSubscription(ctx context.Context, owner, repo string) error {
	return p.subscriberClient.DeleteSubscription(ctx, owner, repo)
}

func (p *SubscriptionProvider) GetSubscriptionsList(ctx context.Context) ([]domain.Subscription, error) {
	return p.subscriberClient.GetSubscriptionsList(ctx)
}

func (p *SubscriptionProvider) GetSubscriptionsInfo(ctx context.Context) ([]domain.Repo, error) {
	return p.processorClient.GetSubscriptionInfo(ctx)
}
