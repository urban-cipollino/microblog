package post

import (
	"github.com/google/wire"

	"micro-blog/internal/handler"
)

var Set = wire.NewSet(
	NewCreatePostService,
	wire.Bind(new(handler.CreatePostService), new(*CreatePostService)),
	NewLikePostService,
	wire.Bind(new(handler.LikePostService), new(*LikePostService)),
	NewPostsService,
	wire.Bind(new(handler.PostsService), new(*PostsService)),
)
