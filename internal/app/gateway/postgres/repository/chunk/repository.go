package chunk

import (
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/postgres/sqlc"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/library/uow"
)

type ChunkRepository struct {
	*sqlc.Queries
}

func NewChunkRepository(db uow.DBTX) *ChunkRepository {
	return &ChunkRepository{
		Queries: sqlc.New(db),
	}
}
