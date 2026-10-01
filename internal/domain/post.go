package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID
	AuthorID  uuid.UUID
	Text      string
	CreatedAt time.Time
}
