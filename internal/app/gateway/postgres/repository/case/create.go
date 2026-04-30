package case_repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/postgres/sqlc"
)

func (r *CaseRepository) Create(ctx context.Context, caseEntity *entity.Case) error {
	const operation = "CaseRepository.Create"

	err := r.CreateCase(ctx, sqlc.CreateCaseParams{
		ID:          caseEntity.ID,
		Name:        caseEntity.Name,
		Description: pgtype.Text{String: caseEntity.Description, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("%s -> %v", operation, err)
	}
	return nil
}
