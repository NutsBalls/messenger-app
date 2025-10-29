package handlers

import (
	"messenger-app/internal/auth/service"
	"messenger-app/internal/auth/store/generated"
	"messenger-app/pkg/hasher"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
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

	passHash, err := hasher.Hash(req.Password)
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
	var req generated.LogInParams

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if req.Email == "" || req.PasswordHash == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "empty request"})
	}

	loginData, err := h.service.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, loginData)
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.Bind(&body); err != nil || body.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userID, err := h.service.ParseToken(c.Request().Context(), body.RefreshToken, "refresh")

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
	}

	user, err := h.service.GetUserByID(c.Request().Context(), userID)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	if user.CryptedRefreshToken == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "no refresh token stored"})
	}

	hashedRefresh, err := hasher.HashToken(body.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad token provided 1"})
	}
	if hashedRefresh != *user.CryptedRefreshToken {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "refresh token mismatch"})
	}

	newTokens, err := h.service.GenerateTokens(c.Request().Context(), user.ID.String())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad token provided 2"})
	}

	cryptedRefresh, err := hasher.HashToken(newTokens.Refresh)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad token provided 3"})
	}

	if err := h.service.UpdateRefreshToken(c.Request().Context(), userID, &cryptedRefresh); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad token provided 4"})
	}

	return c.JSON(http.StatusOK, newTokens)

}
