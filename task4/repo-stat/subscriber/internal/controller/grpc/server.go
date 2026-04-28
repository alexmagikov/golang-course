package grpc

import (
	"context"
	"log/slog"
	subscriberpb "repo-stat/proto/subscriber"
	"repo-stat/subscriber/internal/usecase"
)

type Server struct {
	subscriberpb.UnsafeSubscriberServiceServer
	log                 *slog.Logger
	ping                *usecase.Ping
	subscriptionStorage *usecase.SubscriptionProvider
}

func NewServer(log *slog.Logger, ping *usecase.Ping, subscriptionStorage *usecase.SubscriptionProvider) *Server {
	return &Server{
		log:                 log,
		ping:                ping,
		subscriptionStorage: subscriptionStorage,
	}
}

func (s *Server) Ping(ctx context.Context, _ *subscriberpb.PingRequest) (*subscriberpb.PingResponse, error) {
	s.log.Debug("subscriberp ping request received")

	return &subscriberpb.PingResponse{
		Reply: s.ping.Execute(ctx),
	}, nil
}

func (s *Server) CreateSubscription(ctx context.Context, req *subscriberpb.CreateSubRequest) (*subscriberpb.CreateSubResponse, error) {
	s.log.Debug("subscribers save request received")

	err := s.subscriptionStorage.CreateSubscription(ctx, req.Url)
	if err != nil {
		return nil, err
	}

	return &subscriberpb.CreateSubResponse{}, nil
}

func (s *Server) DeleteSubscription(ctx context.Context, req *subscriberpb.DeleteSubRequest) (*subscriberpb.DeleteSubResponse, error) {
	s.log.Debug("subscribers save request received")

	err := s.subscriptionStorage.DeleteSubscription(ctx, req.Owner, req.Repo)
	if err != nil {
		return nil, err
	}

	return &subscriberpb.DeleteSubResponse{}, nil
}

func (s *Server) GetSubscriptionsList(ctx context.Context, _ *subscriberpb.GetSubListRequest) (*subscriberpb.GetSubListResponse, error) {
	s.log.Debug("subscribers save request received")

	subs, err := s.subscriptionStorage.GetSubscriptionsList(ctx)
	if err != nil {
		return nil, err
	}

	subscriptions := make([]*subscriberpb.Subscription, 0, len(subs))
	for _, sub := range subs {
		subscriptions = append(subscriptions, &subscriberpb.Subscription{
			Id:    sub.ID,
			Owner: sub.Owner,
			Repo:  sub.Repo,
		})
	}

	return &subscriberpb.GetSubListResponse{
		Subscriptions: subscriptions,
	}, nil
}
