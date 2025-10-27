package handlers

import (
	"messenger-app/internal/auth/service"
	"messenger-app/internal/auth/store/generated"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	service *service.AuthService
}

type RegistredRequest struct {
	Name     string
	Password string
	Email    string
}

type LoginRequest struct {
	Email    string
	Password string
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	var req RegistredRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if len([]byte(req.Email)) <= 7 && !strings.Contains(req.Email, "@") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "wrong email"})
	}

	passHash, err := HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	params := generated.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passHash,
	}

	user, err := h.service.SignUp(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Создан пользователь",
		"ID":      user.ID,
		"email":   user.Email,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if req.Email == "" {
		req.Email = c.QueryParam("email")
	}
	if req.Password == "" {
		req.Password = c.QueryParam("password")
	}

	user := generated.LogInParams{
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	token, user2, err := h.service.Login(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"email":        user2.Email,
		"access token": token,
	})
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
