package document

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/gateway/postgres/sqlc"
)

func (r *DocumentRepository) Save(ctx context.Context, doc *entity.Document) error {
	return r.SaveDocument(ctx, sqlc.SaveDocumentParams{
		ID:          doc.ID,
		CaseID:      doc.CaseID,
		FileName:    doc.FileName,
		Path:        doc.Path,
		MimeType:    doc.MimeType,
		SizeBytes:   doc.SizeBytes,
		Status:      int32(doc.Status),
		ChunksCount: int32(doc.ChunksCount),
		LastError: pgtype.Text{
			String: valueOrEmpty(doc.LastError),
			Valid:  doc.LastError != nil,
		},
	})
}

func valueOrEmpty(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
