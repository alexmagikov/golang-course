package grpcController

import (
	"context"
	"repo-stat/processor/internal/usecase"
	"repo-stat/proto/processor"
)

type PingServer struct {
	processor.UnimplementedProcessorServiceServer
	usecase *usecase.Ping
}

func NewPingServer(usecase *usecase.Ping) *PingServer {
	return &PingServer{
		usecase: usecase,
	}
}

func (s *PingServer) Ping(ctx context.Context, request *processor.PingRequest) (*processor.PingResponse, error) {
	status, err := s.usecase.Ping(ctx)
	if err != nil {
		return &processor.PingResponse{Reply: "down"}, err
	}

	return &processor.PingResponse{
		Reply: string(status),
	}, nil
}
