package case_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
)

type caseRepository interface {
	ListAll(ctx context.Context) ([]*entity.Case, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Case, error)
	Create(ctx context.Context, c *entity.Case) error
}
