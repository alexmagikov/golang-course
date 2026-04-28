package usecase

import (
	"context"
	"repo-stat/processor/internal/domain"
)

type SubscriptionInfo struct {
	collectorClient CollectorClient
}

func NewSubscriptionInfo(collectorClient CollectorClient) *SubscriptionInfo {
	return &SubscriptionInfo{collectorClient}
}

func (u *SubscriptionInfo) GetSubscriptionInfo(ctx context.Context) ([]domain.Repo, error) {
	return u.collectorClient.GetSubscriptionInfo(ctx)
}
