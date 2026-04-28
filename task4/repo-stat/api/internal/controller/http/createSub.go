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

// CreateSubscriptionHandler
// @Summary Create sub to repo
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param request body dto.URLRequest true "URL of git repo in JSON"
// @Success 200 {object} string "Success"
// @Failure 400 {string} string "Incorrect request"
// @Failure 404 {string} string "Repo not found"
// @Failure 500 {string} string "Server error"
// @Failure 503 {string} string "Service is unavailable"
// @Router  /api/subscriptions [post]
func CreateSubscriptionHandler(log *slog.Logger, subscriptionProvider *usecase.SubscriptionProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req dto.URLRequest

		errUrl := json.NewDecoder(r.Body).Decode(&req)
		if errUrl != nil {
			log.Error("failed to decode request", "err", errUrl)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to decode request"})
			return
		}

		if req.URL == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "url is required",
			})
			return
		}

		err := subscriptionProvider.CreateSubscription(r.Context(), req.URL)
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
				case codes.AlreadyExists:
					statusCode = http.StatusBadRequest
					errorMsg = "sub already exists"
				}
			} else {
				statusCode = http.StatusBadRequest
				errorMsg = "Invalid url"
			}

			log.Error("create sub failed", "error", err, "url", req.URL)
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
