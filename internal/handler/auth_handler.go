package handler

import (
	"encoding/json"
	"net/http"

	"micro-blog/internal/exception"
	"micro-blog/internal/request"
	"micro-blog/internal/response"
	"micro-blog/internal/service/auth"
)

type AuthHandler struct {
	loginService *auth.LoginService
}

func NewAuthHandler(loginService *auth.LoginService) *AuthHandler {
	return &AuthHandler{
		loginService: loginService,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &exception.BadRequestError{Msg: err.Error()}
	}

	if req.Name == "" || req.Password == "" {
		return &exception.BadRequestError{Msg: "name and password are required"}
	}

	accessToken, err := h.loginService.Login(ctx, req.Name, req.Password)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(response.AuthResponse{AccessToken: accessToken}) //nolint:gosec
}
