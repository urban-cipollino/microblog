package app

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"micro-blog/internal/config"
	"micro-blog/internal/logger"
)

type App struct {
	Port   config.Port
	Router http.Handler
	Logger logger.Logger
	Env    logger.Env
	DB     *pgxpool.Pool
}

func NewApp(
	port config.Port,
	router http.Handler,
	logger logger.Logger,
	env logger.Env,
	db *pgxpool.Pool,
) *App {
	return &App{
		Port:   port,
		Router: router,
		Logger: logger,
		Env:    env,
		DB:     db,
	}
}
