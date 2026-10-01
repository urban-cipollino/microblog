package post

import (
	"context"

	"github.com/google/uuid"

	"micro-blog/internal/domain"
	"micro-blog/internal/dto"
	"micro-blog/internal/logger"
	"micro-blog/internal/repository"
)

type CreatePostService struct {
	userRepo *repository.UserRepository
	postRepo *repository.PostRepository
	logger   logger.Logger
}

func NewCreatePostService(
	userRepo *repository.UserRepository,
	postRepo *repository.PostRepository,
	logger logger.Logger,
) *CreatePostService {
	return &CreatePostService{
		userRepo: userRepo,
		postRepo: postRepo,
		logger:   logger,
	}
}

func (s *CreatePostService) CreatePost(ctx context.Context, authorID uuid.UUID, text string) (*dto.PostDTO, error) {
	author, err := s.userRepo.Get(ctx, authorID)
	if err != nil {
		s.logger.Info("author not found", "author.id", authorID)
		return nil, err
	}

	p := &domain.Post{
		ID:       uuid.New(),
		AuthorID: authorID,
		Text:     text,
	}

	created, err := s.postRepo.Create(ctx, p)
	if err != nil {
		s.logger.Info("failed to create post", "error", err, "author.id", authorID)
		return nil, err
	}

	s.logger.Info("post created", "post.id", created.ID)

	return dto.NewPostDTO(created, author, 0), nil
}
