package post

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"micro-blog/internal/domain"
	"micro-blog/internal/dto"
	"micro-blog/internal/logger"
	"micro-blog/internal/repository"
)

type likeJob struct {
	postID        uuid.UUID
	postCreatedAt time.Time
	userID        uuid.UUID
}

type LikePostService struct {
	db       *pgxpool.Pool
	postRepo *repository.PostRepository
	userRepo *repository.UserRepository
	likeRepo *repository.PostLikeRepository
	logger   logger.Logger
	jobs     chan likeJob
}

func NewLikePostService(
	db *pgxpool.Pool,
	postRepo *repository.PostRepository,
	userRepo *repository.UserRepository,
	likeRepo *repository.PostLikeRepository,
	logger logger.Logger,
) *LikePostService {
	s := &LikePostService{
		db:       db,
		postRepo: postRepo,
		userRepo: userRepo,
		likeRepo: likeRepo,
		logger:   logger,
		jobs:     make(chan likeJob, 64),
	}

	go s.worker()

	return s
}

func (s *LikePostService) worker() {
	for job := range s.jobs {
		ctx := context.Background()

		like := &domain.PostLike{
			PostID:        job.postID,
			PostCreatedAt: job.postCreatedAt,
			UserID:        job.userID,
		}

		created, err := s.likeRepo.Create(ctx, like)
		if err != nil {
			s.logger.Info("failed to save like", "error", err, "post.id", job.postID, "user.id", job.userID)
			continue
		}

		if created == nil {
			s.logger.Info("post is already liked", "post.id", job.postID, "user.id", job.userID)
			continue
		}

		s.logger.Info("like saved", "post.id", job.postID, "user.id", job.userID)
	}
}

func (s *LikePostService) LikePost(ctx context.Context, postID, userID uuid.UUID) (*dto.PostDTO, error) {
	post, err := s.postRepo.Get(ctx, postID)
	if err != nil {
		s.logger.Info("post not found", "post.id", postID)
		return nil, err
	}

	author, err := s.userRepo.Get(ctx, post.AuthorID)
	if err != nil {
		s.logger.Info("author not found", "author.id", post.AuthorID)
		return nil, err
	}

	likes, err := s.likeRepo.CountByPost(ctx, postID, post.CreatedAt)
	if err != nil {
		s.logger.Info("failed to count likes", "error", err, "post.id", postID)
		return nil, err
	}

	s.jobs <- likeJob{
		postID:        post.ID,
		postCreatedAt: post.CreatedAt,
		userID:        userID,
	}

	return dto.NewPostDTO(post, author, likes), nil
}
