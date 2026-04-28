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
		Name:        resp.Repo.Name,
		Description: resp.Repo.Description,
		Stars:       resp.Repo.Stars,
		Forks:       resp.Repo.Forks,
		CreatedAt:   resp.Repo.CreatedAt,
	}, nil
}

func (c *Client) GetSubscriptionInfo(ctx context.Context) ([]domain.Repo, error) {
	resp, err := c.pb.GetSubscriptionsInfo(ctx, &collector.GetSubInfoRequest{})
	if err != nil {
		c.log.Error(err.Error())
		return nil, err
	}

	repos := make([]domain.Repo, 0, len(resp.Repos))
	for _, repo := range resp.Repos {
		repos = append(repos, domain.Repo{
			Name:        repo.Name,
			Description: repo.Description,
			Stars:       repo.Stars,
			Forks:       repo.Forks,
			CreatedAt:   repo.CreatedAt,
		})
	}

	return repos, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
