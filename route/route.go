package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/wire"

	"micro-blog/internal/handler"
	"micro-blog/internal/logger"
	"micro-blog/internal/middleware"
	"micro-blog/internal/service/auth"
)

var Set = wire.NewSet(
	NewHandlers,
	NewRouter,
)

type Handlers struct {
	JWTService  *auth.JWTService
	AuthHandler *handler.AuthHandler
	UserHandler *handler.UserHandler
	PostHandler *handler.PostHandler
	logger      logger.Logger
}

func NewHandlers(
	jwtService *auth.JWTService,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	logger logger.Logger,
) *Handlers {
	return &Handlers{
		JWTService:  jwtService,
		AuthHandler: authHandler,
		UserHandler: userHandler,
		PostHandler: postHandler,
		logger:      logger,
	}
}

func NewRouter(handlers *Handlers) http.Handler {
	r := chi.NewRouter()

	register(r, handlers)

	return r
}

func register(r chi.Router, h *Handlers) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello world!")); err != nil {
			h.logger.Error("failed to write response", "error", err)
		}
	})

	addPost(r, "/api/users", h.UserHandler.CreateUserRequest)
	addPost(r, "/api/login", h.AuthHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(h.JWTService, h.logger))

		addPost(r, "/api/posts", h.PostHandler.CreatePostRequest)
		addGet(r, "/api/posts", h.PostHandler.PostsQuery)
		addPost(r, "/api/posts/{postId}/like", h.PostHandler.LikePostRequest)
	})
}

func addPost(r chi.Router, path string, h middleware.HandleWrapper) {
	r.Method("POST", path, h)
}

func addGet(r chi.Router, path string, h middleware.HandleWrapper) {
	r.Method("GET", path, h)
}
