package usecase

import (
	"context"
	"repo-stat/collector/internal/domain"
)

type RepoProvider struct {
	githubClient GitHubClient
}

func NewRepoProvider(githubClient GitHubClient) *RepoProvider {
	return &RepoProvider{
		githubClient: githubClient,
	}
}

func (u *RepoProvider) GetRepo(ctx context.Context, owner, name string) (domain.Repo, error) {
	return u.githubClient.GetRepo(ctx, owner, name)
}
