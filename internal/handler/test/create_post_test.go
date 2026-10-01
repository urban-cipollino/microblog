package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"micro-blog/internal/domain"
	"micro-blog/internal/dto"
	"micro-blog/internal/handler"
	"micro-blog/internal/middleware"
	"micro-blog/test_help"
)

type mockCreatePostService struct {
	post *dto.PostDTO
	err  error
}

func (m *mockCreatePostService) CreatePost(_ context.Context, _ uuid.UUID, _ string) (*dto.PostDTO, error) {
	return m.post, m.err
}

func newPostRouter(svc handler.CreatePostService, userID uuid.UUID) http.Handler {
	r := chi.NewRouter()
	h := handler.NewPostHandler(svc, nil, nil)
	r.Group(func(r chi.Router) {
		r.Use(test_help.FakeAuth(userID))
		r.Method(http.MethodPost, "/api/posts", middleware.HandleWrapper(h.CreatePostRequest))
	})
	return r
}

func newPostRouterNoAuth(svc handler.CreatePostService) http.Handler {
	r := chi.NewRouter()
	h := handler.NewPostHandler(svc, nil, nil)
	r.Method(http.MethodPost, "/api/posts", middleware.HandleWrapper(h.CreatePostRequest))
	return r
}

func TestCreatePostRequestHandler(t *testing.T) {
	author := &domain.User{ID: uuid.New(), Name: "Jackie Boy"}
	post := &domain.Post{ID: uuid.New(), AuthorID: author.ID, Text: "Ethanol"}

	t.Run("creates post successfully", func(t *testing.T) {
		svc := &mockCreatePostService{
			post: dto.NewPostDTO(post, author, 0),
		}

		body, err := json.Marshal(map[string]string{"text": "Ethanol"})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		newPostRouter(svc, author.ID).ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("returns 401 without auth", func(t *testing.T) {
		svc := &mockCreatePostService{}

		body, err := json.Marshal(map[string]string{"text": "Ethanol"})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		newPostRouterNoAuth(svc).ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("returns 400 on invalid json", func(t *testing.T) {
		svc := &mockCreatePostService{}

		req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBufferString(`invalid`))
		w := httptest.NewRecorder()

		newPostRouter(svc, uuid.New()).ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns 422 on missing text", func(t *testing.T) {
		svc := &mockCreatePostService{}

		body, err := json.Marshal(map[string]string{})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		newPostRouter(svc, uuid.New()).ServeHTTP(w, req)

		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}
