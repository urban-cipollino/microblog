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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"micro-blog/internal/domain"
	"micro-blog/internal/handler"
	"micro-blog/internal/middleware"
	"micro-blog/internal/response"
)

type mockCreateUserService struct {
	user *domain.User
	err  error
}

func (m *mockCreateUserService) CreateUser(_ context.Context, name, _ string) (*domain.User, error) {
	if m.user != nil {
		m.user.Name = name
	}
	return m.user, m.err
}

func newUserRouter(svc handler.CreateUserService) http.Handler {
	r := chi.NewRouter()
	h := handler.NewUserHandler(svc)
	r.Method(http.MethodPost, "/api/users", middleware.HandleWrapper(h.CreateUserRequest))
	return r
}

func TestCreateUserRequestHandler(t *testing.T) {
	t.Run("creates user successfully", func(t *testing.T) {
		svc := &mockCreateUserService{
			user: &domain.User{ID: uuid.New(), Name: "Jackie Boy"},
		}

		body := bytes.NewBufferString(`{"name": "Jackie Boy", "password": "secret123"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/users", body)
		w := httptest.NewRecorder()

		newUserRouter(svc).ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		var u response.UserResponse
		err := json.NewDecoder(w.Body).Decode(&u)
		require.NoError(t, err)

		assert.Equal(t, "Jackie Boy", u.Name)
	})

	t.Run("returns 400 on invalid json", func(t *testing.T) {
		svc := &mockCreateUserService{}

		body := bytes.NewBufferString(`invalid json`)
		req := httptest.NewRequest(http.MethodPost, "/api/users", body)
		w := httptest.NewRecorder()

		newUserRouter(svc).ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns 422 on missing name", func(t *testing.T) {
		svc := &mockCreateUserService{}

		body := bytes.NewBufferString(`{"password": "secret123"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/users", body)
		w := httptest.NewRecorder()

		newUserRouter(svc).ServeHTTP(w, req)

		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}
