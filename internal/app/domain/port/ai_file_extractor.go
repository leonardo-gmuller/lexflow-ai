package port

import "context"

type AIFileExtractor interface {
	Extract(ctx context.Context, mimeType string, fileName string, content []byte) (string, error)
}
