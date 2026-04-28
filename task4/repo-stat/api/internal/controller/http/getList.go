package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

// GetSubscriptionsListHandler GetListOfSubs
// @Summary Get list of current subs
// @Tags Subscriptions
// @Produce json
// @Success 200 {array}  dto.SubscriptionResponse
// @Failure 400 {string} string "Incorrect request"
// @Failure 404 {string} string "Repo not found"
// @Failure 500 {string} string "Server error"
// @Failure 503 {string} string "Service is unavailable"
// @Router  /api/subscriptions [get]
func GetSubscriptionsListHandler(log *slog.Logger, subscriptionProvider *usecase.SubscriptionProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subs, err := subscriptionProvider.GetSubscriptionsList(r.Context())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "service unavailable",
			})
			return
		}

		response := make([]dto.SubscriptionResponse, len(subs))
		for i, s := range subs {
			response[i] = dto.SubscriptionResponse{
				ID:    s.ID,
				Owner: s.Owner,
				Repo:  s.Repo,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}
