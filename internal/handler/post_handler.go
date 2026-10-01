package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"micro-blog/internal/dto"
	"micro-blog/internal/exception"
	"micro-blog/internal/request"
	"micro-blog/internal/response"
	"micro-blog/internal/service/auth"
	"micro-blog/internal/validator"
)

type CreatePostService interface {
	CreatePost(ctx context.Context, authorID uuid.UUID, text string) (*dto.PostDTO, error)
}

type LikePostService interface {
	LikePost(ctx context.Context, postID, userID uuid.UUID) (*dto.PostDTO, error)
}

type PostsService interface {
	Posts(ctx context.Context) ([]*dto.PostDTO, error)
}

type PostHandler struct {
	createPost CreatePostService
	likePost   LikePostService
	posts      PostsService
}

func NewPostHandler(
	createPost CreatePostService,
	likePost LikePostService,
	posts PostsService,
) *PostHandler {
	return &PostHandler{
		createPost: createPost,
		likePost:   likePost,
		posts:      posts,
	}
}

func (ph *PostHandler) CreatePostRequest(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userId, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return &exception.UnauthorizedError{Msg: "user not found"}
	}

	var req request.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &exception.BadRequestError{Msg: err.Error()}
	}

	if err := validator.Check(req); err != nil {
		return err
	}

	created, err := ph.createPost.CreatePost(ctx, userId, req.Text)
	if err != nil {
		return err
	}

	resp := response.NewPostResponse(created.Post, created.Author, created.Likes)

	w.WriteHeader(http.StatusCreated)

	return json.NewEncoder(w).Encode(resp)
}

func (ph *PostHandler) PostsQuery(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	posts, err := ph.posts.Posts(ctx)
	if err != nil {
		return err
	}

	resps := make([]response.PostResponse, 0, len(posts))
	for _, v := range posts {
		resps = append(resps, response.NewPostResponse(v.Post, v.Author, v.Likes))
	}

	return json.NewEncoder(w).Encode(resps)
}

func (ph *PostHandler) LikePostRequest(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	postId, err := uuid.Parse(chi.URLParam(r, "postId"))
	if err != nil {
		return &exception.BadRequestError{Msg: err.Error()}
	}

	userId, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return &exception.UnauthorizedError{Msg: "user not found"}
	}

	post, err := ph.likePost.LikePost(ctx, postId, userId)
	if err != nil {
		return err
	}

	resp := response.NewPostResponse(post.Post, post.Author, post.Likes)

	return json.NewEncoder(w).Encode(resp)
}
