package test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"micro-blog/internal/logger"
	"micro-blog/internal/repository"
	"micro-blog/internal/service/auth"
	"micro-blog/internal/service/post"
	"micro-blog/internal/service/user"
	"micro-blog/test_help"
)

type services struct {
	pool       *pgxpool.Pool
	createUser *user.CreateUserService
	createPost *post.CreatePostService
	likePost   *post.LikePostService
	posts      *post.PostsService
}

func setup(t testing.TB) *services {
	t.Helper()

	pool := test_help.SetupPostgres(t)
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	log := logger.New(logger.Test)

	userRepo := repository.NewUserRepository(pool, qb)
	postRepo := repository.NewPostRepository(pool, qb)
	likeRepo := repository.NewPostLikeRepository(pool, qb)
	passwordSvc := auth.NewPasswordService()

	return &services{
		pool:       pool,
		createUser: user.NewCreateUserService(log, userRepo, passwordSvc),
		createPost: post.NewCreatePostService(userRepo, postRepo, log),
		likePost:   post.NewLikePostService(pool, postRepo, userRepo, likeRepo, log),
		posts:      post.NewPostsService(postRepo),
	}
}
