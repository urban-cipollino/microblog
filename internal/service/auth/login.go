package auth

import (
	"context"

	"micro-blog/internal/exception"
	"micro-blog/internal/logger"
	"micro-blog/internal/repository"
)

type LoginService struct {
	logger          logger.Logger
	userRepo        *repository.UserRepository
	passwordService *PasswordService
	jwtService      *JWTService
}

func NewLoginService(
	logger logger.Logger,
	userRepo *repository.UserRepository,
	passwordService *PasswordService,
	jwtService *JWTService,
) *LoginService {
	return &LoginService{
		userRepo:        userRepo,
		logger:          logger,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (s LoginService) Login(ctx context.Context, name, password string) (string, error) {
	user, err := s.userRepo.GetByName(ctx, name)
	if err != nil {
		return "", err
	}

	if !s.passwordService.Compare(user.Password, password) {
		return "", &exception.UnauthorizedError{Msg: "invalid credentials"}
	}

	accessToken, err := s.jwtService.Generate(user.ID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
