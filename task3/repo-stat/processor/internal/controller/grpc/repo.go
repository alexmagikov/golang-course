package grpcController

import (
	"context"
	"repo-stat/processor/internal/usecase"
	"repo-stat/proto/processor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RepoServer struct {
	processor.UnimplementedProcessorServiceServer
	usecase *usecase.Repo
}

func NewRepoServer(usecase *usecase.Repo) *RepoServer {
	return &RepoServer{
		usecase: usecase,
	}
}

func (s *RepoServer) GetRepoInfo(ctx context.Context, req *processor.GetRepoRequest) (*processor.GetRepoResponse, error) {
	if req.GetUrl() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty repo url")
	}

	resp, err := s.usecase.GetRepoInfo(ctx, req.GetUrl())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &processor.GetRepoResponse{
		Name:        resp.Name,
		Description: resp.Description,
		Stars:       resp.Stars,
		Forks:       resp.Forks,
		CreatedAt:   resp.CreatedAt,
	}, nil
}
