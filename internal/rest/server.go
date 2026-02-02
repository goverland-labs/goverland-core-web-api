package rest

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/goverland-labs/goverland-core-web-api/internal/response"

	apihandlers "github.com/goverland-labs/goverland-core-web-api/internal/rest/handlers"

	"github.com/goverland-labs/goverland-core-web-api/internal/config"
	"github.com/goverland-labs/goverland-core-web-api/pkg/middleware"
)

func NewRestServer(cfg config.REST, apiHandlers []apihandlers.APIHandler) *http.Server {
	handler := mux.NewRouter()
	handler.Use(
		middleware.Panic,
		middleware.Prometheus,
		middleware.ResponseFormatter,
	)

	baseV1Router := handler.PathPrefix("/v1").Subrouter()
	baseV1Router.Use(middleware.Timeout(cfg.HandleTimeout))

	baseV2Router := handler.PathPrefix("/v2").Subrouter()
	baseV2Router.Use(middleware.Timeout(cfg.HandleTimeout))

	for _, h := range apiHandlers {
		h.EnrichRoutes(baseV1Router, baseV2Router)
	}

	return &http.Server{
		Addr:         cfg.Listen,
		Handler:      configureCorsHandler(handler),
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	}
}

func configureCorsHandler(router *mux.Router) http.Handler {
	handlerMethods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodDelete})
	handlerCredentials := handlers.AllowCredentials()
	handlerAllowedHeaders := handlers.AllowedHeaders([]string{
		"Content-Type",
		"Authorization",
	})
	handlerExposedHeaders := handlers.ExposedHeaders([]string{
		response.HeaderTotalCount,
		response.HeaderCurrentOffset,
		response.HeaderLimit,
	})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	return handlers.CORS(handlerMethods, handlerCredentials, handlerAllowedHeaders, handlerExposedHeaders, allowedOrigins)(router)
}
