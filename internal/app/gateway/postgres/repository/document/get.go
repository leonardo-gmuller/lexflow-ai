package document

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
)

func (r *DocumentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
	row, err := r.FindDocumentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.Document{
		ID:          row.ID,
		CaseID:      row.CaseID,
		FileName:    row.FileName,
		Path:        row.Path,
		MimeType:    row.MimeType,
		SizeBytes:   row.SizeBytes,
		Status:      entity.DocumentStatus(row.Status),
		ChunksCount: int(row.ChunksCount),
		LastError:   textPtr(row.LastError),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}, nil
}

func textPtr(v pgtype.Text) *string {
	if !v.Valid {
		return nil
	}
	return &v.String
}
