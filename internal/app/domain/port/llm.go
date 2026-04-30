package port

import (
	"context"

	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/dto"
)

type LLM interface {
	GenerateResponse(ctx context.Context, prompt string) (dto.LLMResponse, error)
}
