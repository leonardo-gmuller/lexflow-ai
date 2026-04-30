package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/dto"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/entity"
)

func (uc *ProcessDocument) Execute(ctx context.Context, job dto.ProcessDocumentJob) error {

	// 1. buscar documento
	doc, err := uc.docRepo.FindByID(ctx, job.DocumentID)
	if err != nil {
		return err
	}

	// 2. marcar como processing
	_ = uc.docRepo.UpdateStatus(ctx, doc.ID, entity.DocumentStatusProcessing, nil)

	// 3. ler arquivo
	content, err := uc.storage.Read(ctx, doc.Path)
	if err != nil {
		return uc.fail(ctx, doc.ID, err)
	}

	// 4. extrair texto
	text, err := uc.extractor.Extract(ctx, doc.MimeType, doc.FileName, content)
	if err != nil {
		return uc.fail(ctx, doc.ID, err)
	}

	// 5. chunking
	chunksText := uc.chunker.Split(text)

	if len(chunksText) == 0 {
		return uc.fail(ctx, doc.ID, fmt.Errorf("no chunks generated"))
	}

	// 6. embeddings
	_, err = uc.embedder.Generate(ctx, chunksText)
	if err != nil {
		return uc.fail(ctx, doc.ID, err)
	}

	// 7. montar chunks
	var chunks []*entity.DocumentChunk

	for i, text := range chunksText {
		chunks = append(chunks, &entity.DocumentChunk{
			ID:         uuid.New(),
			DocumentID: doc.ID,
			CaseID:     doc.CaseID,
			Content:    text,
			ChunkIndex: i,
		})
	}

	// 8. salvar chunks
	if err := uc.chunkRepo.SaveMany(ctx, chunks); err != nil {
		return uc.fail(ctx, doc.ID, err)
	}

	// 9. atualizar documento
	if err := uc.docRepo.UpdateChunksCount(ctx, doc.ID, len(chunks)); err != nil {
		return err
	}

	if err := uc.docRepo.UpdateStatus(ctx, doc.ID, entity.DocumentStatusProcessed, nil); err != nil {
		return err
	}

	return nil
}
