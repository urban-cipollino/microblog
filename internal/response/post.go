package response

import (
	"time"

	"github.com/google/uuid"

	"micro-blog/internal/domain"
)

type PostResponse struct {
	ID        uuid.UUID    `json:"id"`
	Author    UserResponse `json:"author"`
	Text      string       `json:"text"`
	Likes     int          `json:"likes"`
	CreatedAt time.Time    `json:"created_at"`
}

func NewPostResponse(p *domain.Post, author *domain.User, likes int) PostResponse {
	return PostResponse{
		ID:        p.ID,
		Author:    NewUserResponse(author),
		Text:      p.Text,
		Likes:     likes,
		CreatedAt: p.CreatedAt,
	}
}
