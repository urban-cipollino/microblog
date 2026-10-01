package user

import (
	"context"

	"github.com/google/uuid"

	"micro-blog/internal/domain"
	"micro-blog/internal/exception"
	"micro-blog/internal/logger"
	"micro-blog/internal/repository"
	"micro-blog/internal/service/auth"
)

type CreateUserService struct {
	logger          logger.Logger
	userRepo        *repository.UserRepository
	passwordService *auth.PasswordService
}

func NewCreateUserService(
	logger logger.Logger,
	userRepo *repository.UserRepository,
	passwordService *auth.PasswordService,
) *CreateUserService {
	return &CreateUserService{
		logger:          logger,
		userRepo:        userRepo,
		passwordService: passwordService,
	}
}

func (s *CreateUserService) CreateUser(ctx context.Context, name, password string) (*domain.User, error) {
	if s.userRepo.IsNameTaken(ctx, name) {
		s.logger.Info("user creation failed: name is taken", "user.name", name)

		return nil, &exception.ConflictError{Msg: "name is taken"}
	}

	hash, err := s.passwordService.Hash(password)
	if err != nil {
		s.logger.Info("failed to hash password", "error", err, "user.password")
		return nil, err
	}

	u := &domain.User{ID: uuid.New(), Name: name, Password: hash}
	if _, err := s.userRepo.Create(ctx, u); err != nil {
		s.logger.Info("failed to create user", "error", err, "user.name", name)

		return nil, err
	}

	s.logger.Info("user created", "user.id", u.ID, "user.name", u.Name)

	return u, nil
}
