package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/leonardo-gmuller/lexflow-ai/internal/app"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/config"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/api"
	open_ai "github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/openai"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/postgres"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/redis"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/storage"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/library/logger"
	"golang.org/x/sync/errgroup"
)

var Version = "dev"

func main() {
	ctx := context.Background()
	logger.Init()

	// 1) Config
	cfg, err := config.New(Version)
	if err != nil {
		log.Fatalf("failed to load configurations: %v", err)
	}

	// 2) Postgres
	postgresClient, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to start postgres: %v", err)
	}
	defer postgresClient.Close()

	// 3) Redis
	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		log.Fatalf("failed to start redis: %v", err)
	}

	//MCP
	// mcp := tools.NewToolEngine(repository.NewBotConfigRepository(postgresClient.Pool))

	// 5) OpenAI Client (LLM)
	openAI, err := open_ai.NewOpenAI(cfg.OpenAI)
	if err != nil {
		log.Fatalf("failed to start OpenAI client: %v", err)
	}

	storage := storage.NewLocalFileStorage(cfg.Storage.BasePath)

	// 7) Application (contém PG, Redis, LLM, senders, BatchStore, etc.)
	appl := app.New(
		ctx,
		cfg,
		postgresClient,
		redisClient,
		openAI,
		openAI,  // AIExtractor - implementar depois
		openAI,  // Embedder - implementar depois
		storage, // Storage - implementar depois
	)

	// 10) Graceful Shutdown
	stopCtx, stop := signal.NotifyContext(ctx,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	defer stop()

	group, groupCtx := errgroup.WithContext(stopCtx)

	// 9) Server HTTP (API)
	server := &http.Server{
		Addr: cfg.Server.Address,
		BaseContext: func(_ net.Listener) context.Context {
			return stopCtx
		},
		Handler:      api.New(cfg, appl).Handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// goroutine do servidor HTTP
	group.Go(func() error {
		log.Printf("starting api server on %s", cfg.Server.Address)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	// goroutine de shutdown gracioso
	group.Go(func() error {
		<-groupCtx.Done()

		log.Printf("stopping api; interrupt signal received")

		timeoutCtx, cancel := context.WithTimeout(context.Background(), cfg.App.GracefulShutdownTimeout)
		defer cancel()

		var errs error

		if err := server.Shutdown(timeoutCtx); err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to stop server: %w", err))
		}

		return errs
	})

	if err := group.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("api exit reason: %v", err)
	}
}
