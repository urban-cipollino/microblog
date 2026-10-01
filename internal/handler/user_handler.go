package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"micro-blog/internal/domain"
	"micro-blog/internal/exception"
	"micro-blog/internal/request"
	"micro-blog/internal/response"
	"micro-blog/internal/validator"
)

type CreateUserService interface {
	CreateUser(ctx context.Context, name, password string) (*domain.User, error)
}

type UserHandler struct {
	createUser CreateUserService
}

func NewUserHandler(createUser CreateUserService) *UserHandler {
	return &UserHandler{
		createUser: createUser,
	}
}

func (uh *UserHandler) CreateUserRequest(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req request.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &exception.BadRequestError{Msg: err.Error()}
	}

	if err := validator.Check(req); err != nil {
		return err
	}

	created, err := uh.createUser.CreateUser(ctx, req.Name, req.Password)
	if err != nil {
		return err
	}

	resp := response.NewUserResponse(created)

	w.WriteHeader(http.StatusCreated)

	return json.NewEncoder(w).Encode(resp)
}
