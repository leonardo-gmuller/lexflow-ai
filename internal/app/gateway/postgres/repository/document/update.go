package document

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/postgres/sqlc"
)

func (r *DocumentRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status entity.DocumentStatus,
	lastError *string,
) error {
	return r.UpdateDocumentStatus(ctx, sqlc.UpdateDocumentStatusParams{
		ID:     id,
		Status: int32(status),
		LastError: pgtype.Text{
			String: valueOrEmpty(lastError),
			Valid:  lastError != nil,
		},
	})
}

func (r *DocumentRepository) UpdateChunksCount(ctx context.Context, id uuid.UUID, count int) error {
	return r.UpdateDocumentChunksCount(ctx, sqlc.UpdateDocumentChunksCountParams{
		ID:          id,
		ChunksCount: int32(count),
	})
}
