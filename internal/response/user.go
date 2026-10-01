package response

import (
	"github.com/google/uuid"

	"micro-blog/internal/domain"
)

type UserResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func NewUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:   u.ID,
		Name: u.Name,
	}
}
