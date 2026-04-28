package usecase

import (
	"context"
	"log/slog"
	"repo-stat/subscriber/internal/domain"
	"repo-stat/subscriber/internal/repository/postgres"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GitHubClient interface {
	CheckRepo(ctx context.Context, owner, repo string) (bool, error)
}

type SubscriptionProvider struct {
	log          *slog.Logger
	GithubClient GitHubClient
	queries      *postgres.Queries
}

func NewSubscriptionProvider(log *slog.Logger, githubClient GitHubClient, queries *postgres.Queries) *SubscriptionProvider {
	return &SubscriptionProvider{
		log:          log,
		GithubClient: githubClient,
		queries:      queries,
	}
}

func parseURL(url string) (owner, repo string, err error) {
	url = strings.TrimPrefix(url, "https://github.com/")
	url = strings.TrimPrefix(url, "github.com/")
	parts := strings.Split(url, "/")
	if len(parts) != 2 {
		return "", "", status.Error(codes.InvalidArgument, "invalid url")
	}

	return parts[0], parts[1], nil
}

func (u *SubscriptionProvider) CreateSubscription(ctx context.Context, url string) error {
	owner, repo, err := parseURL(url)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	u.log.Debug("checking repo on github", "owner", owner, "repo", repo)

	exists, err := u.GithubClient.CheckRepo(ctx, owner, repo)
	if err != nil {
		u.log.Error("github check failed", "error", err)
		return err
	}
	if !exists {
		return status.Error(codes.NotFound, "repository not found")

	}

	_, errSub := u.queries.GetSubscriptionsByRepo(ctx, postgres.GetSubscriptionsByRepoParams{
		Owner: owner,
		Repo:  repo,
	})

	if errSub == nil {
		return status.Error(codes.AlreadyExists, "sub already exists")
	}

	_, err = u.queries.CreateSubscription(ctx, postgres.CreateSubscriptionParams{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		u.log.Error("create subscription failed", "error", err)
		return err
	}

	u.log.Info("sub saved", "owner", owner, "repo", repo)
	return nil
}

func (u *SubscriptionProvider) DeleteSubscription(ctx context.Context, owner, repo string) error {
	err := u.queries.DeleteSubscriptions(ctx, postgres.DeleteSubscriptionsParams{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return err
	}

	u.log.Info("sub deleted", "owner", owner, "repo", repo)
	return nil
}

func (u *SubscriptionProvider) GetSubscriptionsList(ctx context.Context) ([]domain.Subscription, error) {
	subs, err := u.queries.ListSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Subscription, 0, len(subs))
	for _, sub := range subs {
		result = append(result, domain.Subscription{
			ID:    sub.ID,
			Owner: sub.Owner,
			Repo:  sub.Repo,
		})
	}

	return result, nil
}
