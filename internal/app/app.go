package app

import (
	"context"

	"github.com/leonardo-gmuller/lexflow-ai/internal/app/config"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/port"
	case_usecase "github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/usecase/case"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/usecase/document"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/postgres"
	redisGt "github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/redis"

	case_repository "github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/postgres/repository/case"
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

func New(ctx context.Context,
	cfg config.Config,
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
		redisGt.NewRedisDocumentQueue(a.Queue),
	)
}

func (a *App) NewCaseUseCase(dbtx postgres.DBTX) case_usecase.CaseUseCaseInterface {
	return case_usecase.NewCaseUseCase(
		case_repository.NewCaseRepository(dbtx),
	)
}
