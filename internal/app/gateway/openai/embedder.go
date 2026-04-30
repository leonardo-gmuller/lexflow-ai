package open_ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/leonardo-gmuller/lexflow-ai/internal/app/library/util"
	openai "github.com/openai/openai-go/v3"
)

const defaultEmbeddingModel = "text-embedding-3-small"

func (c *OpenAIClient) EmbedText(ctx context.Context, text string) ([]float32, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("text vazio")
	}

	model := defaultEmbeddingModel

	resp, err := c.cli.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: openai.EmbeddingModel(model),
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("resposta de embedding vazia")
	}

	return util.Float64SliceToFloat32(resp.Data[0].Embedding), nil
}

func (c *OpenAIClient) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	cleaned := make([]string, 0, len(texts))
	for _, t := range texts {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		cleaned = append(cleaned, t)
	}

	if len(cleaned) == 0 {
		return nil, fmt.Errorf("nenhum texto válido para embedding")
	}

	model := defaultEmbeddingModel

	resp, err := c.cli.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: openai.EmbeddingModel(model),
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: cleaned,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar embeddings em lote: %w", err)
	}

	if len(resp.Data) != len(cleaned) {
		return nil, fmt.Errorf("quantidade de embeddings inconsistente: textos=%d embeddings=%d", len(cleaned), len(resp.Data))
	}

	out := make([][]float32, 0, len(resp.Data))
	for _, item := range resp.Data {
		out = append(out, util.Float64SliceToFloat32(item.Embedding))
	}

	return out, nil
}
