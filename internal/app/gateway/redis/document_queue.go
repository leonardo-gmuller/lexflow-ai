package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/dto"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/usecase/document"
	redisv9 "github.com/redis/go-redis/v9"
)

const redisDocumentScheduledJobsKey = "lexflow:document:scheduled_jobs"

type DocumentQueue struct {
	client ClientInterface
}

func NewRedisDocumentQueue(client ClientInterface) document.Queue {
	return &DocumentQueue{
		client: client,
	}
}

func (s *DocumentQueue) Publish(
	ctx context.Context,
	documentID uuid.UUID,
	delay time.Duration,
) error {
	job := dto.ProcessDocumentJob{
		DocumentID: documentID,
		Attempt:    0,
	}

	raw, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal knowledge job: %w", err)
	}

	execAt := time.Now().Add(delay).Unix()

	if err := s.client.ZAdd(ctx, redisDocumentScheduledJobsKey, redisv9.Z{
		Score:  float64(execAt),
		Member: string(raw),
	}); err != nil {
		return fmt.Errorf("redis zadd scheduled document: %w", err)
	}

	return nil
}

func (s *DocumentQueue) FetchDueDocuments(
	ctx context.Context,
	now time.Time,
	limit int64,
) ([]dto.ProcessDocumentJob, error) {
	maxScore := now.Unix()

	values, err := s.client.ZRangeByScore(
		ctx,
		redisDocumentScheduledJobsKey,
		"-inf",
		strconv.FormatInt(maxScore, 10),
		0,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("redis zrangebyscore knowledge: %w", err)
	}
	if len(values) == 0 {
		return nil, nil
	}

	if err := s.client.ZRem(ctx, redisDocumentScheduledJobsKey, values...); err != nil {
		// best effort, igual ao WhatsApp
	}

	out := make([]dto.ProcessDocumentJob, 0, len(values))
	for _, raw := range values {
		var job dto.ProcessDocumentJob
		if err := json.Unmarshal([]byte(raw), &job); err != nil {
			continue
		}
		out = append(out, job)
	}

	return out, nil
}

func (s *DocumentQueue) RetryDocument(
	ctx context.Context,
	job dto.ProcessDocumentJob,
	delay time.Duration,
) error {
	job.Attempt++

	raw, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal retry knowledge job: %w", err)
	}

	execAt := time.Now().Add(delay).Unix()

	if err := s.client.ZAdd(ctx, redisDocumentScheduledJobsKey, redisv9.Z{
		Score:  float64(execAt),
		Member: string(raw),
	}); err != nil {
		return fmt.Errorf("redis zadd retry document: %w", err)
	}

	return nil
}
