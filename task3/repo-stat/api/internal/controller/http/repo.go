package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewRepoHandler
// @Summary Get info about repo by url
// @Tags Repositories
// @Accept json
// @Produce json
// @Param url query string true "url of github repo"
// @Success 200 {object} dto.RepoResponse "Success"
// @Failure 400 {string} string "Incorrect request"
// @Failure 404 {string} string "Repo not found"
// @Failure 500 {string} string "Server error"
// @Failure 503 {string} string "Service is unavailable"
// @Router  /api/repositories/info [get]
func NewRepoHandler(log *slog.Logger, repo *usecase.RepoUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("url is required")
			return
		}

		repo, err := repo.GetRepoInfo(r.Context(), url)
		if err != nil {
			statusCode := http.StatusInternalServerError
			errorMsg := "Internal Server Error"

			if st, ok := status.FromError(err); ok {
				switch st.Code() {
				case codes.NotFound:
					statusCode = http.StatusNotFound
					errorMsg = "repo not found"
				case codes.InvalidArgument:
					statusCode = http.StatusBadRequest
					errorMsg = "invalid url"
				case codes.Unavailable:
					statusCode = http.StatusServiceUnavailable
					errorMsg = "service unavailable"
				}
			}

			log.Error("get repo info failed", "error", err, "url", url)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(map[string]string{
				"error": errorMsg,
			})

			return
		}

		response := dto.FromDomainRepo(repo)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("encode response failed", "error", err)
		}

	}
}
