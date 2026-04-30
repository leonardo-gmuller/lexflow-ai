package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/config"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/api/handler"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/api/middleware"
)

type API struct {
	Handler http.Handler
	cfg     config.Config
	app     *app.App
}

func New(cfg config.Config, a *app.App) *API {
	api := &API{
		cfg: cfg,
		app: a,
	}

	api.setupRouter()

	return api
}

func (api *API) setupRouter() {
	router := chi.NewRouter()

	if api.cfg.Development {
		router.Use(middleware.Logger)
	}

	router.Use(
		middleware.CORS,
		middleware.CleanPath,
		middleware.StripSlashes,
		middleware.Recoverer,
		middleware.RealIP,
		middleware.HTTPErrorLogger,
	)

	api.registerRoutes(router)

	api.Handler = router
}

func BasicHandler() http.Handler {
	router := chi.NewMux()
	return router
}

func (api *API) registerRoutes(router *chi.Mux) {
	// health check público
	handler.RegisterBasicRoutes(router, api.cfg, api.app)

	router.Route("/api/v1/lexflow", func(publicRouter chi.Router) {
		handler.RegisterExternalRoutes(
			publicRouter,
			api.cfg,
			api.app,
			api.app.Queue,
		)
	})
}
