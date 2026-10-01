package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"micro-blog/internal/domain"
	"micro-blog/internal/dto"
	"micro-blog/internal/handler"
	"micro-blog/internal/middleware"
	"micro-blog/internal/response"
	"micro-blog/test_help"
)

type mockLikePostService struct {
	post *dto.PostDTO
	err  error
}

func (m *mockLikePostService) LikePost(_ context.Context, _, _ uuid.UUID) (*dto.PostDTO, error) {
	return m.post, m.err
}

func newLikeRouter(svc handler.LikePostService, userID uuid.UUID) http.Handler {
	r := chi.NewRouter()
	h := handler.NewPostHandler(nil, svc, nil)
	r.Group(func(r chi.Router) {
		r.Use(test_help.FakeAuth(userID))
		r.Method(http.MethodPost, "/api/posts/{postId}/like", middleware.HandleWrapper(h.LikePostRequest))
	})
	return r
}

func newLikeRouterNoAuth(svc handler.LikePostService) http.Handler {
	r := chi.NewRouter()
	h := handler.NewPostHandler(nil, svc, nil)
	r.Method(http.MethodPost, "/api/posts/{postId}/like", middleware.HandleWrapper(h.LikePostRequest))
	return r
}

func TestLikePostRequestHandler(t *testing.T) {
	t.Run("likes post", func(t *testing.T) {
		author := &domain.User{ID: uuid.New(), Name: "Jackie Boy"}
		post := &domain.Post{ID: uuid.New(), AuthorID: author.ID, Text: "test"}
		userID := uuid.New()

		svc := &mockLikePostService{
			post: dto.NewPostDTO(post, author, 1),
		}

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/posts/"+post.ID.String()+"/like",
			nil,
		)
		w := httptest.NewRecorder()

		newLikeRouter(svc, userID).ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var p response.PostResponse
		err := json.NewDecoder(w.Body).Decode(&p)
		require.NoError(t, err)

		assert.Equal(t, post.ID, p.ID)
		assert.Equal(t, 1, p.Likes)
	})

	t.Run("returns 401 without auth", func(t *testing.T) {
		svc := &mockLikePostService{}

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/posts/"+uuid.New().String()+"/like",
			nil,
		)
		w := httptest.NewRecorder()

		newLikeRouterNoAuth(svc).ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// BenchmarkLikePostRequestHandler-32 126212 8858 ns/op 7481 B/op 30 allocs/op
func BenchmarkLikePostRequestHandler(b *testing.B) {
	author := &domain.User{ID: uuid.New(), Name: "Jackie Boy"}
	post := &domain.Post{ID: uuid.New(), AuthorID: author.ID, Text: "test"}
	userID := uuid.New()

	svc := &mockLikePostService{
		post: dto.NewPostDTO(post, author, 1),
	}
	router := newLikeRouter(svc, userID)
	url := "/api/posts/" + post.ID.String() + "/like"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, url, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(b, http.StatusOK, w.Code)
	}
}
