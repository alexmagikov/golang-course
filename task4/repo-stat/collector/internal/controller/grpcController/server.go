package grpcController

import (
	"context"
	"repo-stat/collector/internal/usecase"
	"repo-stat/proto/collector"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CollectorServer struct {
	collector.UnimplementedCollectorServiceServer
	uc *usecase.CollectorProvider
}

func NewCollectorServer(uc *usecase.CollectorProvider) *CollectorServer {
	return &CollectorServer{
		uc: uc,
	}
}

func (s *CollectorServer) GetRepo(ctx context.Context, req *collector.GetRepoRequest) (*collector.GetRepoResponse, error) {
	if req.GetOwner() == "" || req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "owner, name, required")
	}

	repo, err := s.uc.GetRepo(ctx, req.GetOwner(), req.GetName())
	if err != nil {
		return nil, err
	}

	return &collector.GetRepoResponse{
		Repo: &collector.Repo{
			Name:        repo.Name,
			Description: repo.Description,
			Stars:       repo.Stars,
			Forks:       repo.Forks,
			CreatedAt:   repo.CreatedAt,
		},
	}, nil
}

func (s *CollectorServer) GetSubscriptionsInfo(ctx context.Context, req *collector.GetSubInfoRequest) (*collector.GetSubInfoResponse, error) {
	subs, err := s.uc.GetSubscriptionsInfo(ctx)
	if err != nil {
		return nil, err
	}

	repos := make([]*collector.Repo, 0, len(subs))
	for _, repo := range subs {
		repos = append(repos, &collector.Repo{
			Name:        repo.Name,
			Description: repo.Description,
			Stars:       repo.Stars,
			Forks:       repo.Forks,
			CreatedAt:   repo.CreatedAt,
		})
	}

	return &collector.GetSubInfoResponse{
		Repos: repos,
	}, nil
}
