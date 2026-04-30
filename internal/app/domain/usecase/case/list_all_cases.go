package case_usecase

import (
	"context"

	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/dto"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
)

type ListAllCasesInput struct {
	dto.Pagination
	Search string
}

type CasesOrderBy string

const (
	OrderByName      CasesOrderBy = "name"
	OrderByCreatedAt CasesOrderBy = "created_at"
)

type ListAllCasesOutput struct {
	Cases      []entity.Case `json:"cases"`
	TotalItems int           `json:"total_items"`
}

func (uc *CaseUseCase) ListAllCases(ctx context.Context, in ListAllCasesInput) (*ListAllCasesOutput, error) {
	cases, totalItems, err := uc.repo.ListAll(
		ctx,
		in.Search,
		(string)(CasesOrderBy(in.SortBy)),
		in.SortType,
		in.Page,
		in.Offset(),
	)

	if err != nil {
		return nil, err
	}

	return &ListAllCasesOutput{
		Cases:      cases,
		TotalItems: totalItems,
	}, nil
}
