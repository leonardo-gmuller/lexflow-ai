package port

import "context"

type Storage interface {
	Save(ctx context.Context, path string, content []byte) error
	Read(ctx context.Context, path string) ([]byte, error)
}
