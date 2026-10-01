//go:build wireinject
// +build wireinject

package app

import (
	"micro-blog/internal/config"
	"micro-blog/internal/database"
	"micro-blog/internal/handler"
	"micro-blog/internal/logger"
	"micro-blog/internal/repository"
	authservice "micro-blog/internal/service/auth"
	postservice "micro-blog/internal/service/post"
	userservice "micro-blog/internal/service/user"
	"micro-blog/route"

	"github.com/google/wire"
)

func InitializeApp() *App {
	wire.Build(
		config.Set,
		database.Set,
		logger.Set,
		authservice.Set,
		repository.Set,
		userservice.Set,
		postservice.Set,
		handler.Set,
		route.Set,
		NewApp,
	)

	return nil
}
