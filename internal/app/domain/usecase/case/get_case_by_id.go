package case_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/erring"
)

func (uc *CaseUseCase) GetCaseByID(ctx context.Context, id uuid.UUID) (*entity.Case, error) {
	c, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, erring.ErrCaseNotFound
	}
	return c, nil
}
