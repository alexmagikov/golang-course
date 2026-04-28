package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	httpClient *http.Client
	token      string
}

func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		token: token,
	}
}

func (c *Client) CheckRepo(ctx context.Context, owner, repo string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "GetRepoInfo-App")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:

	case http.StatusNotFound:
		return false, status.Error(codes.NotFound, "repo not found")

	case http.StatusForbidden:
		return false, status.Error(codes.ResourceExhausted, "access forbidden")

	case http.StatusUnauthorized:
		return false, status.Error(codes.Unauthenticated, "unauthorized: invalid token")

	default:
		return false, status.Error(codes.Internal, "unexpected status code")
	}

	return true, nil
}
