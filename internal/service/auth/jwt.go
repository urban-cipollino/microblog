package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"micro-blog/internal/config"
	"micro-blog/internal/exception"
	"micro-blog/internal/logger"
)

type JWTService struct {
	logger logger.Logger
	secret []byte
	ttl    time.Duration
}

func NewJWTService(
	logger logger.Logger,
	secret config.JwtSecret,
	ttl config.JwtTTL,
) *JWTService {
	return &JWTService{
		logger: logger,
		secret: []byte(secret),
		ttl:    time.Duration(ttl),
	}
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type contextKey string

const UserIDKey contextKey = "user_id"

func (s *JWTService) Generate(userID uuid.UUID) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secret)
}

func (s *JWTService) Parse(tokenString string) (uuid.UUID, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			s.logger.Warn("jwt: unexpected signing method", "method", token.Method.Alg())
			return nil, &exception.UnauthorizedError{Msg: "unexpected signing method"}
		}
		return s.secret, nil
	})
	if err != nil {
		s.logger.Warn("jwt: parse failed", "error", err)
		return uuid.Nil, &exception.UnauthorizedError{Msg: "invalid token"}
	}

	if !token.Valid {
		s.logger.Warn("jwt: token invalid")
		return uuid.Nil, &exception.UnauthorizedError{Msg: "invalid token"}
	}

	return claims.UserID, nil
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)

	return userID, ok
}
