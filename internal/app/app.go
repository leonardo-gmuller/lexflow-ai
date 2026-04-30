package app

import (
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/config"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/port"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/usecase/document"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/postgres"
	redisGt "github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/redis"

	documentRepo "github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/postgres/repository/document"
)

type App struct {
	Config config.Config

	DB    *postgres.Client
	Queue redisGt.ClientInterface
	LLM   port.LLM

	AIExtractor port.AIFileExtractor
	Embedder    port.Embedder
	Storage     port.Storage
}

func New(cfg config.Config,
	db *postgres.Client,
	queue redisGt.ClientInterface,
	llm port.LLM,
	aiExtractor port.AIFileExtractor,
	embedder port.Embedder,
	storage port.Storage,
) *App {
	return &App{
		Config:      cfg,
		DB:          db,
		Queue:       queue,
		LLM:         llm,
		AIExtractor: aiExtractor,
		Embedder:    embedder,
		Storage:     storage,
	}
}

func (a *App) NewDocumentUsecase(dbtx postgres.DBTX) document.DocumentUsecaseInterface {
	return document.NewDocumentUsecase(
		documentRepo.NewDocumentRepository(dbtx),
		a.Storage,
		redisGt.NewRedisDocumentQueue(a.Queue),
	)
}

func (a *App) NewProcessDocumentUsecase(dbtx postgres.DBTX) document.ProcessDocumentInterface {
	return document.NewProcessDocument(
		documentRepo.NewDocumentRepository(dbtx),
		nil, // chunkRepo - implementar depois
		a.Storage,
		a.AIExtractor,
		nil, // chunker - implementar depois
		a.Embedder,
	)
}
