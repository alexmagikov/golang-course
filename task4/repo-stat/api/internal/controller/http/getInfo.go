package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

// GetSubscriptionInfoHandler
// @Summary Get info about repo by url
// @Tags Subscriptions
// @Produce json
// @Success 200 {array} dto.RepoResponse "Success"
// @Failure 400 {string} string "Incorrect request"
// @Failure 404 {string} string "Repo not found"
// @Failure 500 {string} string "Server error"
// @Failure 503 {string} string "Service is unavailable"
// @Router  /api/subscriptions/info [get]
func GetSubscriptionInfoHandler(log *slog.Logger, subscriptionProvider *usecase.SubscriptionProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repos, err := subscriptionProvider.GetSubscriptionsInfo(r.Context())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "service unavailable",
			})
			return
		}

		response := make([]dto.RepoInfoResponse, len(repos))
		for i, repo := range repos {
			response[i] = dto.FromDomainRepoInfo(repo)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("encode response failed", "error", err)
		}

	}
}
