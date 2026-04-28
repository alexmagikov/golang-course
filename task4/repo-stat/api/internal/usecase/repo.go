package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type RepoUseCase struct {
	processorClient ProcessorClient
}

func NewRepo(processorClient ProcessorClient) *RepoUseCase {
	return &RepoUseCase{
		processorClient: processorClient,
	}
}

func (u *RepoUseCase) GetRepoInfo(ctx context.Context, url string) (domain.Repo, error) {
	return u.processorClient.GetRepo(ctx, url)
}
