package document

import (
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/library/uow"

	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/postgres/sqlc"
)

type DocumentRepository struct {
	*sqlc.Queries
}

func NewDocumentRepository(db uow.DBTX) *DocumentRepository {
	return &DocumentRepository{Queries: sqlc.New(db)}
}
