package case_repository

import (
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/postgres/sqlc"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/library/uow"
)

type CaseRepository struct {
	*sqlc.Queries
}

func NewCaseRepository(db uow.DBTX) *CaseRepository {
	return &CaseRepository{Queries: sqlc.New(db)}
}
