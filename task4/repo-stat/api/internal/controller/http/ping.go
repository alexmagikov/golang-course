package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

// NewPingHandler
// @Summary Ping services
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} dto.PingResponse "All services are available"
// @Failure 503 {object} dto.PingResponse "One or all services are not available"
// @Router /api/ping [get]
func NewPingHandler(log *slog.Logger, ping *usecase.PingUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := ping.Execute(r.Context())

		response := dto.FromDomainPingResult(result)

		statusCode := http.StatusOK
		if result.Status == "degraded" {
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write ping response", "error", err)
		}
	}
}
