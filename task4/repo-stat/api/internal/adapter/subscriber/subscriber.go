package subscriber

import (
	"context"
	"log/slog"
	"repo-stat/api/internal/domain"

	subscriberpb "repo-stat/proto/subscriber"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   subscriberpb.SubscriberServiceClient
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
		pb:   subscriberpb.NewSubscriberServiceClient(conn),
	}, nil
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.pb.Ping(ctx, &subscriberpb.PingRequest{})
	if err != nil {
		c.log.Error("subscriber ping failed", "error", err)
		return domain.PingStatusDown
	}

	return domain.PingStatusUp
}

func (c *Client) CreateSubscription(ctx context.Context, url string) error {
	_, err := c.pb.CreateSubscription(ctx, &subscriberpb.CreateSubRequest{Url: url})
	if err != nil {
		c.log.Error("subscriber create subscription failed", "error", err)
		return err
	}

	return nil
}

func (c *Client) DeleteSubscription(ctx context.Context, owner, repo string) error {
	_, err := c.pb.DeleteSubscription(ctx, &subscriberpb.DeleteSubRequest{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		c.log.Error("subscriber create subscription failed", "error", err)
		return err
	}

	return nil
}

func (c *Client) GetSubscriptionsList(ctx context.Context) ([]domain.Subscription, error) {
	subs, err := c.pb.GetSubscriptionsList(ctx, &subscriberpb.GetSubListRequest{})
	if err != nil {
		c.log.Error("subscriber get subscriptions list failed", "error", err)
		return nil, err
	}

	subscriptions := make([]domain.Subscription, 0, len(subs.Subscriptions))
	for _, sub := range subs.Subscriptions {
		subscriptions = append(subscriptions, domain.Subscription{
			ID:    sub.Id,
			Owner: sub.Owner,
			Repo:  sub.Repo,
		})
	}

	return subscriptions, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
