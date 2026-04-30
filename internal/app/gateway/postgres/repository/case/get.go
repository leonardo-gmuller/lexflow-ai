package case_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/postgres/sqlc"
)

func (r *CaseRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Case, error) {
	const operation = "CaseRepository.GetByID"

	caseData, err := r.FindCaseByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s -> %v", operation, err)
	}

	return &entity.Case{
		ID:          caseData.ID,
		Name:        caseData.Name,
		Description: caseData.Description.String,
		CreatedAt:   caseData.CreatedAt.Time,
		UpdatedAt:   caseData.UpdatedAt.Time,
	}, nil
}

func (r *CaseRepository) ListAll(ctx context.Context, search, sortBy, sortOrder string, page, offset int) ([]entity.Case, int, error) {
	const operation = "CaseRepository.ListAll"

	rows, err := r.ListCases(ctx, sqlc.ListCasesParams{
		SearchTerm: pgtype.Text{String: search, Valid: true},
		SortBy:     sortBy,
		SortOrder:  sortOrder,
		SqlLimit:   int32(page),
		SqlOffset:  int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("%s -> %v", operation, err)
	}

	var totalCount int
	if len(rows) > 0 {
		totalCount = int(rows[0].TotalCount)
	}

	cases := make([]entity.Case, 0, len(rows))
	for _, row := range rows {
		cases = append(cases, entity.Case{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description.String,
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   row.UpdatedAt.Time,
		})
	}

	return cases, totalCount, nil
}
