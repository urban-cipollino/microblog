package dto

import (
	"micro-blog/internal/domain"
)

type PostDTO struct {
	Post   *domain.Post
	Author *domain.User
	Likes  int
}

func NewPostDTO(post *domain.Post, author *domain.User, likes int) *PostDTO {
	return &PostDTO{Post: post, Author: author, Likes: likes}
}
