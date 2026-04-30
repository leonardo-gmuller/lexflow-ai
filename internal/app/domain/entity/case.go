package entity

import (
	"time"

	"github.com/google/uuid"
)

type Case struct {
	ID   uuid.UUID
	Name string

	Description string

	CreatedAt time.Time
	UpdatedAt time.Time
}
