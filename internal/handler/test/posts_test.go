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
)

type mockPostsService struct {
	posts []*dto.PostDTO
	err   error
}

func (m *mockPostsService) Posts(_ context.Context) ([]*dto.PostDTO, error) {
	return m.posts, m.err
}

func newPostsRouter(svc handler.PostsService) http.Handler {
	r := chi.NewRouter()
	h := handler.NewPostHandler(nil, nil, svc)
	r.Method(http.MethodGet, "/api/posts", middleware.HandleWrapper(h.PostsQuery))
	return r
}

func TestPostsQueryHandler(t *testing.T) {
	t.Run("returns posts", func(t *testing.T) {
		author := &domain.User{ID: uuid.New(), Name: "Jackie Boy"}
		post := &domain.Post{ID: uuid.New(), AuthorID: author.ID, Text: "test"}
		svc := &mockPostsService{
			posts: []*dto.PostDTO{dto.NewPostDTO(post, author, 0)},
		}

		req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
		w := httptest.NewRecorder()

		newPostsRouter(svc).ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var posts []response.PostResponse
		err := json.NewDecoder(w.Body).Decode(&posts)
		require.NoError(t, err)

		require.Len(t, posts, 1)
		assert.Equal(t, post.ID, posts[0].ID)
		assert.Equal(t, "test", posts[0].Text)
	})

	t.Run("empty", func(t *testing.T) {
		svc := &mockPostsService{posts: []*dto.PostDTO{}}

		req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
		w := httptest.NewRecorder()

		newPostsRouter(svc).ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var posts []response.PostResponse
		err := json.NewDecoder(w.Body).Decode(&posts)
		require.NoError(t, err)

		assert.Empty(t, posts)
	})
}

// BenchmarkPostsQueryHandler-32 35030 33951 ns/op 41554 B/op 223 allocs/op
func BenchmarkPostsQueryHandler(b *testing.B) {
	posts := make([]*dto.PostDTO, 100)
	for i := range posts {
		author := &domain.User{ID: uuid.New(), Name: "Jackie Boy"}
		post := &domain.Post{ID: uuid.New(), AuthorID: author.ID, Text: "test"}
		posts[i] = dto.NewPostDTO(post, author, 0)
	}

	svc := &mockPostsService{posts: posts}
	router := newPostsRouter(svc)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(b, http.StatusOK, w.Code)
	}
}
