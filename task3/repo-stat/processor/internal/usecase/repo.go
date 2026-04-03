package usecase

import (
	"context"
	"net/url"
	"repo-stat/processor/internal/domain"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Repo struct {
	collectorClient CollectorClient
}

func NewRepo(collectorClient CollectorClient) *Repo {
	return &Repo{
		collectorClient: collectorClient,
	}
}

func ParseGitHubURL(gihubURL string) (owner, name string, err error) {
	parsedURL, err := url.Parse(gihubURL)
	if err != nil {
		return "", "", err
	}

	path := strings.Trim(parsedURL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", status.Error(codes.InvalidArgument, path)
	}

	return parts[0], parts[1], nil
}

func (u *Repo) GetRepoInfo(ctx context.Context, githubURL string) (domain.Repo, error) {
	owner, name, err := ParseGitHubURL(githubURL)
	if err != nil {
		return domain.Repo{}, err
	}

	return u.collectorClient.GetRepo(ctx, owner, name)
}
