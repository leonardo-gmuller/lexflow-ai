package case_usecase

import (
	"context"

	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
)

func (uc *CaseUseCase) ListAllCases(ctx context.Context) ([]*entity.Case, error) {
	return uc.repo.ListAll(ctx)
}
