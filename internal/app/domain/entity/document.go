package entity

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID       uuid.UUID
	CaseID   uuid.UUID
	FileName string
	Path     string

	MimeType    string
	SizeBytes   int64
	ChunksCount int

	Status DocumentStatus

	LastError *string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type DocumentStatus int

const (
	DocumentStatusPending DocumentStatus = iota
	DocumentStatusProcessing
	DocumentStatusProcessed
	DocumentStatusFailed
)
