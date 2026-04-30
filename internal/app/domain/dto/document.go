package dto

import "github.com/google/uuid"

type ProcessDocumentJob struct {
	DocumentID uuid.UUID `json:"document_id"`
	Attempt    int       `json:"attempt"`
}
