package domain

import (
	"time"

	"github.com/google/uuid"
)

type PostLike struct {
	PostID        uuid.UUID
	PostCreatedAt time.Time
	UserID        uuid.UUID
	CreatedAt     time.Time
}
