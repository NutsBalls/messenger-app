package handlers

import (
	"context"
	"messenger-app/internal/auth/service"
	"messenger-app/internal/auth/store/generated"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

func (h *AuthHandler) CreateUser(c echo.Context) error {
	var req generated.CreateUserParams

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := h.service.CreateUser(context.Background(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Создан пользователь",
		"ID":      user.ID,
		"email":   user.Email,
	})
}
