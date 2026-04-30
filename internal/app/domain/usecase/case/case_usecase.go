package case_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
)

type CaseUseCase struct {
	repo caseRepository
}

type CaseUseCaseInterface interface {
	ListAllCases(ctx context.Context, in ListAllCasesInput) (*ListAllCasesOutput, error)
	GetCaseByID(ctx context.Context, id uuid.UUID) (*entity.Case, error)
	CreateCase(ctx context.Context, name string, description string) (*entity.Case, error)
}

func NewCaseUseCase(repo caseRepository) CaseUseCaseInterface {
	return &CaseUseCase{
		repo: repo,
	}
}
