package collector

import (
	"context"
	"log/slog"
	"repo-stat/processor/internal/domain"
	"repo-stat/proto/collector"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   collector.CollectorServiceClient
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
		pb:   collector.NewCollectorServiceClient(conn),
	}, nil
}

func (c *Client) GetRepo(ctx context.Context, owner, name string) (domain.Repo, error) {
	resp, err := c.pb.GetRepo(ctx, &collector.GetRepoRequest{
		Owner: owner,
		Name:  name,
	})
	if err != nil {
		c.log.Error(err.Error())
		return domain.Repo{}, err
	}

	return domain.Repo{
		Name:        resp.Name,
		Description: resp.Description,
		Stars:       resp.Stars,
		Forks:       resp.Forks,
		CreatedAt:   resp.CreatedAt,
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
