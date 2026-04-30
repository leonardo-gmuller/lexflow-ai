package document

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/dto"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
)

type documentRepository interface {
	Save(ctx context.Context, document *entity.Document) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Document, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.DocumentStatus, lastError *string) error
	UpdateChunksCount(ctx context.Context, id uuid.UUID, chunksCount int) error
}

type caseRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Case, error)
}

type storage interface {
	Save(ctx context.Context, path string, content []byte) error
	Read(ctx context.Context, path string) ([]byte, error)
}

type Queue interface {
	FetchDueDocuments(ctx context.Context, now time.Time, limit int64) ([]dto.ProcessDocumentJob, error)
	Publish(ctx context.Context, documentID uuid.UUID, delay time.Duration) error
}

type chunkRepository interface {
	SaveMany(ctx context.Context, chunks []*entity.DocumentChunk) error
}

type textExtractor interface {
	Extract(ctx context.Context, mimeType string, fileName string, content []byte) (string, error)
}

type chunker interface {
	Split(text string) []string
}
