package subscriber

import (
	"context"
	"log/slog"
	"repo-stat/collector/internal/domain"
	"repo-stat/proto/subscriber"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   subscriber.SubscriberServiceClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   subscriber.NewSubscriberServiceClient(conn),
	}, nil
}

func (c *Client) GetSubscriptionsList(ctx context.Context) ([]domain.Subscription, error) {
	resp, err := c.pb.GetSubscriptionsList(ctx, &subscriber.GetSubListRequest{})
	if err != nil {
		c.log.Error("failed to get subscriptions", "error", err)
		return nil, err
	}

	subscriptions := make([]domain.Subscription, 0, len(resp.Subscriptions))
	for _, s := range resp.Subscriptions {
		subscriptions = append(subscriptions, domain.Subscription{
			ID:    s.Id,
			Owner: s.Owner,
			Repo:  s.Repo,
		})
	}

	return subscriptions, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
