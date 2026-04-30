package port

import "context"

type Embedder interface {
	Generate(ctx context.Context, texts []string) ([][]float32, error)
}
