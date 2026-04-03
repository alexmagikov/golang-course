package usecase

import (
	"context"
	"repo-stat/collector/internal/domain"
)

type GitHubClient interface {
	GetRepo(ctx context.Context, owner, name string) (domain.Repo, error)
}
