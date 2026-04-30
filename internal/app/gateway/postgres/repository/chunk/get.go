package chunk

import (
	"context"

	"github.com/google/uuid"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
)

func (r *ChunkRepository) ListByDocumentID(ctx context.Context, documentID uuid.UUID) ([]*entity.DocumentChunk, error) {
	rows, err := r.ListChunksByDocumentID(ctx, documentID)
	if err != nil {
		return nil, err
	}

	chunks := make([]*entity.DocumentChunk, 0, len(rows))

	for _, row := range rows {
		chunks = append(chunks, &entity.DocumentChunk{
			ID:         row.ID,
			DocumentID: row.DocumentID,
			CaseID:     row.CaseID,
			Content:    row.Content,
			ChunkIndex: int(row.ChunkIndex),
			TokenCount: int(row.TokenCount),
			CreatedAt:  row.CreatedAt.Time,
		})
	}

	return chunks, nil
}
