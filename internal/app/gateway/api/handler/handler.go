package handler

import (
	"context"
	"time"

	"github.com/cep21/circuit/v3"
	"github.com/cep21/circuit/v3/closers/hystrix"
	"github.com/go-chi/chi/v5"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/config"
)

// cache interface continua igual
type cache interface {
	Exists(ctx context.Context, key string) (bool, error)
	Get(ctx context.Context, key string, objByRef any) error
	Set(ctx context.Context, key string, obj any, ttl time.Duration) error
}

// Agora o Handler guarda o *app.App em vez de um UseCase monolítico
type Handler struct {
	circuitManager *circuit.Manager
	cfg            config.Config
	app            *app.App
	cache          cache
}

func New(cfg config.Config, a *app.App, cache cache) Handler {
	hystrixFactory := hystrix.Factory{
		ConfigureOpener: hystrix.ConfigureOpener{
			ErrorThresholdPercentage: int64(cfg.CircuitBreaker.ErrorPercentThreshold),
			RequestVolumeThreshold:   int64(cfg.CircuitBreaker.RequestVolumeThreshold),
		},
		ConfigureCloser: hystrix.ConfigureCloser{
			SleepWindow: cfg.CircuitBreaker.SleepWindow,
		},
	}

	defaultFactory := func(_ string) circuit.Config {
		return circuit.Config{
			Execution: circuit.ExecutionConfig{
				MaxConcurrentRequests: int64(cfg.CircuitBreaker.MaxConcurrentRequests),
				Timeout:               cfg.CircuitBreaker.Timeout,
			},
		}
	}

	circuitManager := &circuit.Manager{
		DefaultCircuitProperties: []circuit.CommandPropertiesConstructor{
			defaultFactory,
			hystrixFactory.Configure,
		},
	}

	return Handler{
		circuitManager: circuitManager,
		cfg:            cfg,
		app:            a,
		cache:          cache,
	}
}

func RegisterBasicRoutes(router chi.Router, cfg config.Config, a *app.App) {
	h := New(cfg, a, nil)
	h.healthcheckSetup(router)
}

func RegisterExternalRoutes(
	router chi.Router,
	cfg config.Config,
	a *app.App,
	cache cache,
) {
	h := New(cfg, a, cache)

	h.caseSetupRoutes(router)

}
