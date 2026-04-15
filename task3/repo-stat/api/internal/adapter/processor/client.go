package processor

import (
	"context"
	"log/slog"
	"repo-stat/api/internal/domain"
	"repo-stat/proto/processor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   processor.ProcessorServiceClient
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
		pb:   processor.NewProcessorServiceClient(conn),
	}, nil
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.pb.Ping(ctx, &processor.PingRequest{})
	if err != nil {
		c.log.Error("processor ping failed", "error", err)
		return domain.PingStatusDown
	}

	return domain.PingStatusUp
}

func (c *Client) GetRepo(ctx context.Context, url string) (domain.Repo, error) {
	resp, err := c.pb.GetRepo(ctx, &processor.GetRepoRequest{Url: url})
	if err != nil {
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
	return c.conn.Close()
}
