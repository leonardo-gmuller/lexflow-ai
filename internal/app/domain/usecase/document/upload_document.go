package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/erring"
)

type UploadDocumentInput struct {
	CaseID   string
	FileName string
	Content  []byte
}

func (uc *DocumentUsecase) UploadDocument(ctx context.Context, input UploadDocumentInput) (*entity.Document, error) {
	docID := uuid.New()
	caseID, err := uuid.Parse(input.CaseID)
	if err != nil {
		return nil, fmt.Errorf("invalid case ID: %w", err)
	}
	path := uc.getDocumentPath(caseID, docID, input.FileName)

	caseEntity, err := uc.caseRepo.GetByID(ctx, caseID)
	if err != nil {
		return nil, erring.ErrCaseNotFound
	}

	//1. Save file in storage
	if err := uc.storage.Save(ctx, path, input.Content); err != nil {
		return nil, fmt.Errorf("failed to save document in storage: %w", err)
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
		return nil, fmt.Errorf("failed to save document metadata: %w", err)
	}

	//3. Publish message to queue for processing
	if err := uc.queue.Publish(ctx, docID, 0); err != nil {
		return nil, fmt.Errorf("failed to publish document upload event: %w", err)
	}

	return document, nil
}
