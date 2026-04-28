package grpcController

import (
	"context"
	"repo-stat/processor/internal/usecase"
	"repo-stat/proto/processor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProcessorServer struct {
	processor.UnimplementedProcessorServiceServer
	ping             *usecase.Ping
	repo             *usecase.Repo
	subscriptionInfo *usecase.SubscriptionInfo
}

func NewProcessorServer(ping *usecase.Ping, repo *usecase.Repo, subscriptionInfo *usecase.SubscriptionInfo) *ProcessorServer {
	return &ProcessorServer{
		ping:             ping,
		repo:             repo,
		subscriptionInfo: subscriptionInfo,
	}
}

func (s *ProcessorServer) Ping(ctx context.Context, req *processor.PingRequest) (*processor.PingResponse, error) {
	status, err := s.ping.Ping(ctx)
	if err != nil {
		return &processor.PingResponse{Reply: "down"}, err
	}

	return &processor.PingResponse{
		Reply: string(status),
	}, nil
}

func (s *ProcessorServer) GetRepo(ctx context.Context, req *processor.GetRepoRequest) (*processor.GetRepoResponse, error) {
	if req.GetUrl() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty repo url")
	}

	repo, err := s.repo.GetRepoInfo(ctx, req.GetUrl())
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, st.Err()
		}

		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	return &processor.GetRepoResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreatedAt:   repo.CreatedAt,
	}, nil
}

func (s *ProcessorServer) GetSubscriptionInfo(ctx context.Context, req *processor.GetSubInfoRequest) (*processor.GetSubInfoResponse, error) {
	subs, err := s.subscriptionInfo.GetSubscriptionInfo(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	repos := make([]*processor.Repo, 0, len(subs))
	for _, repo := range subs {
		repos = append(repos, &processor.Repo{
			Name:        repo.Name,
			Description: repo.Description,
			Stars:       repo.Stars,
			Forks:       repo.Forks,
			CreatedAt:   repo.CreatedAt,
		})
	}

	return &processor.GetSubInfoResponse{
		Repos: repos,
	}, nil
}
