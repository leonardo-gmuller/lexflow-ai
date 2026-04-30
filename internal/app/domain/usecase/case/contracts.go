package case_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
)

type caseRepository interface {
	ListAll(ctx context.Context, search, sortBy, sortOrder string, page, offset int) ([]entity.Case, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Case, error)
	Create(ctx context.Context, c *entity.Case) error
}
