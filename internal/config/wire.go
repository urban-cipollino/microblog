package config

import (
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/google/wire"

	"micro-blog/internal/logger"
)

var Set = wire.NewSet(
	NewEnv,
	NewPort,
	NewDatabaseURL,
	NewJwtTtl,
	NewJwtSecret,
)

func NewEnv() logger.Env {
	env := os.Getenv("APP_ENV")
	if env == "" {
		return logger.Dev
	}

	return logger.Env(env)
}

type Port string

func NewPort() Port {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	return Port(port)
}

type DatabaseURL string

func NewDatabaseURL() DatabaseURL {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   name,
	}

	q := u.Query()
	q.Set("sslmode", sslmode)
	u.RawQuery = q.Encode()

	return DatabaseURL(u.String())
}

type JwtTTL time.Duration

func NewJwtTtl() JwtTTL {
	ttl := os.Getenv("JWT_TTL")
	seconds, err := strconv.Atoi(ttl)
	if err != nil {
		seconds = 86400
	}

	return JwtTTL(time.Duration(seconds) * time.Second)
}

type JwtSecret string

func NewJwtSecret() JwtSecret {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is not set")
	}

	return JwtSecret(secret)
}
