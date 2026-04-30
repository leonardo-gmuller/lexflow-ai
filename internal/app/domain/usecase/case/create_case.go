package case_usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
)

func (uc *CaseUseCase) CreateCase(ctx context.Context, name, description string) (*entity.Case, error) {
	c := &entity.Case{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.repo.Create(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}
