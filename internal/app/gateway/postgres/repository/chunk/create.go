package chunk

import (
	"context"

	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/postgres/sqlc"
	"github.com/pgvector/pgvector-go"
)

func (r *ChunkRepository) SaveMany(ctx context.Context, chunks []*entity.DocumentChunk) error {
	for _, chunk := range chunks {
		err := r.SaveDocumentChunk(ctx, sqlc.SaveDocumentChunkParams{
			ID:         chunk.ID,
			DocumentID: chunk.DocumentID,
			CaseID:     chunk.CaseID,
			Content:    chunk.Content,
			ChunkIndex: int32(chunk.ChunkIndex),
			TokenCount: int32(chunk.TokenCount),
			Embedding:  pgvector.Vector{},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
