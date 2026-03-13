package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type RepoInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}

func CreatePath(name string, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s", name, repo)
}

// GetRepoInfo / Get general info about GitHub repo by url.
func GetRepoInfo(url string) (*RepoInfo, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("User-Agent", "GetRepoInfo-App")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:

	case http.StatusNotFound:

		return nil, fmt.Errorf("repo not found: %s", url)

	case http.StatusForbidden:
		return nil, fmt.Errorf("access is not accepted: %s", url)

	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var repo RepoInfo

	if err := json.Unmarshal(body, &repo); err != nil {
		return nil, err
	}

	return &repo, nil
}
