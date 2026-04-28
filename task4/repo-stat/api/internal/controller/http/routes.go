package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/usecase"

	httpSwagger "github.com/swaggo/http-swagger"
)

func AddRoutes(
	mux *http.ServeMux,
	log *slog.Logger,
	ping *usecase.PingUseCase,
	repo *usecase.RepoUseCase,
	subscriptionsProvider *usecase.SubscriptionProvider,
) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping))
	mux.Handle("GET /api/repositories/info", NewRepoHandler(log, repo))

	mux.Handle("POST /api/subscriptions", CreateSubscriptionHandler(log, subscriptionsProvider))
	mux.Handle("DELETE /api/subscriptions/{owner}/{repo}", DeleteSubscriptionHandler(log, subscriptionsProvider))
	mux.Handle("GET /api/subscriptions", GetSubscriptionsListHandler(log, subscriptionsProvider))
	mux.Handle("GET /api/subscriptions/info", GetSubscriptionInfoHandler(log, subscriptionsProvider))

	mux.Handle("GET /swagger/doc.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.json")
	}))

	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))
}
