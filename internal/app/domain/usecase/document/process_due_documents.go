package document

import (
	"context"
	"sync"
	"time"

	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/dto"
)

func (uc *ProcessDocument) ProcessDueDocuments(ctx context.Context, now time.Time, limit int64, maxConcurrency int) error {
	jobs, err := uc.queue.FetchDueDocuments(ctx, now, limit)
	if err != nil {
		return err
	}

	if len(jobs) == 0 {
		return nil
	}

	if maxConcurrency <= 0 {
		maxConcurrency = 1
	}

	// Process documents concurrently
	sem := make(chan struct{}, maxConcurrency)
	errCh := make(chan error, len(jobs))

	var wg sync.WaitGroup

	for _, job := range jobs {
		sem <- struct{}{}
		wg.Add(1)
		go func(job dto.ProcessDocumentJob) {
			defer func() { <-sem }()
			defer wg.Done()
			if err := uc.Execute(ctx, job); err != nil {
				errCh <- err
			}
		}(job)
	}

	// Wait for all goroutines to finish
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}

	close(errCh)
	if len(errCh) > 0 {
		return <-errCh
	}

	return nil
}
