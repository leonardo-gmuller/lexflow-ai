package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalFileStorage struct {
	basePath string
}

func NewLocalFileStorage(basePath string) *LocalFileStorage {
	return &LocalFileStorage{
		basePath: basePath,
	}
}

func (s *LocalFileStorage) fullPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(s.basePath, path)
}

func (s *LocalFileStorage) Save(
	ctx context.Context,
	path string,
	file io.Reader,
) error {

	full := s.fullPath(path)

	dir := filepath.Dir(full)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create storage directory: %w", err)
	}

	out, err := os.Create(full)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (s *LocalFileStorage) Open(
	ctx context.Context,
	path string,
) (io.ReadCloser, error) {
	full := s.fullPath(path)

	file, err := os.Open(full)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %w", full, err)
	}

	return file, nil
}

func (s *LocalFileStorage) Delete(
	ctx context.Context,
	path string,
) error {

	full := s.fullPath(path)

	err := os.Remove(full)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("delete file %s: %w", full, err)
	}

	return nil
}
