package dto

import "repo-stat/api/internal/domain"

type RepoResponse struct {
	FullName    string `json:"full_name"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int64  `json:"stars"`
	Forks       int64  `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

func FromDomainRepo(repo domain.Repo, fullName string) RepoResponse {
	return RepoResponse{
		FullName:    fullName,
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreatedAt:   repo.CreatedAt,
	}
}

type RepoInfoResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int64  `json:"stars"`
	Forks       int64  `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

func FromDomainRepoInfo(repoInfo domain.Repo) RepoInfoResponse {
	return RepoInfoResponse{
		Name:        repoInfo.Name,
		Description: repoInfo.Description,
		Stars:       repoInfo.Stars,
		Forks:       repoInfo.Forks,
		CreatedAt:   repoInfo.CreatedAt,
	}
}
