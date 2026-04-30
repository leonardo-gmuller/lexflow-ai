package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/erring"
)

type UploadDOcumentInput struct {
	CaseID   uuid.UUID
	FileName string
	Content  []byte
}

func (uc *DocumentUsecase) UploadDocument(ctx context.Context, input UploadDOcumentInput) error {
	docID := uuid.New()
	path := uc.getDocumentPath(input.CaseID, docID, input.FileName)

	caseEntity, err := uc.caseRepo.GetByID(ctx, input.CaseID)
	if err != nil {
		return erring.ErrCaseNotFound
	}

	//1. Save file in storage
	if err := uc.storage.Save(ctx, path, input.Content); err != nil {
		return fmt.Errorf("failed to save document in storage: %w", err)
	}

	//2. Save document metadata in database
	document := &entity.Document{
		ID:       docID,
		CaseID:   caseEntity.ID,
		FileName: input.FileName,
		Path:     path,
		Status:   entity.DocumentStatusPending,
	}

	if err := uc.documentRepo.Save(ctx, document); err != nil {
		return fmt.Errorf("failed to save document metadata: %w", err)
	}

	//3. Publish message to queue for processing
	if err := uc.queue.Publish(ctx, docID, 0); err != nil {
		return fmt.Errorf("failed to publish document upload event: %w", err)
	}

	return nil
}
