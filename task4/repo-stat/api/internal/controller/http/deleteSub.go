package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/usecase"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteSubscriptionHandler
// @Summary Delete sub to repo
// @Tags Subscriptions
// @Produce json
// @Param owner path string true "Repo owner"
// @Param repo path string true "Repo name"
// @Success 200 {object} string "Deleted"
// @Failure 400 {string} string "Incorrect request"
// @Failure 404 {string} string "Repo not found"
// @Failure 500 {string} string "Server error"
// @Failure 503 {string} string "Service is unavailable"
// @Router  /api/subscriptions/{owner}/{repo} [delete]
func DeleteSubscriptionHandler(log *slog.Logger, subscriptionProvider *usecase.SubscriptionProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := strings.TrimPrefix(r.URL.Path, "/api/subscriptions/")
		parts := strings.Split(path, "/")
		if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "owner and repo are required",
			})
			return
		}

		owner := parts[0]
		repo := parts[1]

		err := subscriptionProvider.DeleteSubscription(r.Context(), owner, repo)
		if err != nil {
			statusCode := http.StatusInternalServerError
			errorMsg := "Internal Server Error"

			if st, ok := status.FromError(err); ok {
				switch st.Code() {
				case codes.NotFound:
					statusCode = http.StatusNotFound
					errorMsg = "sub not found"
				case codes.InvalidArgument:
					statusCode = http.StatusBadRequest
					errorMsg = "invalid arg"
				case codes.Unavailable:
					statusCode = http.StatusServiceUnavailable
					errorMsg = "service unavailable"
				}
			} else {
				statusCode = http.StatusBadRequest
				errorMsg = "Invalid"
			}

			log.Error("delete sub failed", "error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": errorMsg,
			})

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
