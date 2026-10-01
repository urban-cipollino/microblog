package post

import (
	"context"

	"micro-blog/internal/dto"
	"micro-blog/internal/repository"
)

type PostsService struct {
	postRepo *repository.PostRepository
}

func NewPostsService(postRepo *repository.PostRepository) *PostsService {
	return &PostsService{
		postRepo: postRepo,
	}
}

func (s *PostsService) Posts(ctx context.Context) ([]*dto.PostDTO, error) {
	return s.postRepo.FindAll(ctx)
}
