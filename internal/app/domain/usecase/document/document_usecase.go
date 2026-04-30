package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/dto"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/port"
)

type DocumentUsecase struct {
	documentRepo documentRepository
	caseRepo     caseRepository
	storage      storage
	queue        Queue
}

type ProcessDocument struct {
	docRepo   documentRepository
	chunkRepo chunkRepository
	storage   storage

	queue Queue

	extractor textExtractor
	chunker   chunker
	embedder  port.Embedder
}

type DocumentUsecaseInterface interface {
	UploadDocument(ctx context.Context, input UploadDocumentInput) (*entity.Document, error)
}

type ProcessDocumentInterface interface {
	Execute(ctx context.Context, job dto.ProcessDocumentJob) error
}

func NewDocumentUsecase(documentRepo documentRepository, storage storage, queue Queue) DocumentUsecaseInterface {
	return &DocumentUsecase{
		documentRepo: documentRepo,
		storage:      storage,
		queue:        queue,
	}
}

func NewProcessDocument(
	docRepo documentRepository,
	chunkRepo chunkRepository,
	storage storage,
	extractor textExtractor,
	chunker chunker,
	embedder port.Embedder,
	queue Queue,
) ProcessDocumentInterface {
	return &ProcessDocument{
		docRepo:   docRepo,
		chunkRepo: chunkRepo,
		storage:   storage,
		extractor: extractor,
		chunker:   chunker,
		embedder:  embedder,
		queue:     queue,
	}
}

func (uc *DocumentUsecase) getDocumentPath(caseID, documentID uuid.UUID, fileName string) string {
	return fmt.Sprintf("cases/%s/documents/%s/%s", caseID.String(), documentID.String(), fileName)
}

func (uc *ProcessDocument) fail(ctx context.Context, docID uuid.UUID, err error) error {
	msg := err.Error()
	_ = uc.docRepo.UpdateStatus(ctx, docID, entity.DocumentStatusFailed, &msg)
	return err
}
