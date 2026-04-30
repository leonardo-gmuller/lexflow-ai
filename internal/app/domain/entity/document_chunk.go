package entity

import (
	"time"

	"github.com/google/uuid"
)

type DocumentChunk struct {
	ID         uuid.UUID
	DocumentID uuid.UUID
	CaseID     uuid.UUID

	Content string

	ChunkIndex int
	TokenCount int

	CreatedAt time.Time
}
