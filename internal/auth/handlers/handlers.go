package handlers

import (
	"messenger-app/internal/auth/service"
	"messenger-app/internal/auth/store/generated"
	"net/http"
	"strings"
	"time"

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

	accessToken, refreshToken, user2, err := h.service.Login(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"email":         user2.Email,
		"access token":  accessToken,
		"refresh token": refreshToken,
	})
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.Bind(&body); err != nil || body.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	check, err := h.service.CheckRefreshToken(c.Request().Context(), body.RefreshToken)
	if err != nil || !check {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
	}

	userID, err := h.service.ParseToken(c.Request().Context(), body.RefreshToken, "refresh")

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
	}

	user, err := h.service.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "user not found"})
	}

	newAccessToken, err := h.service.GenerateAccessToken(c.Request().Context(), user.ID.String(), time.Duration(time.Minute*30))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create new access token"})
	}

	newRefreshToken, err := h.service.GenerateRefreshToken(c.Request().Context(), userID.String(), time.Duration(time.Hour*24))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create new refresh token"})
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	if err := h.service.UpdateRefreshToken(c.Request().Context(), user.ID.Bytes, newRefreshToken, expiresAt); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to add new refresh token in dab"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access token":  newAccessToken,
		"refresh token": newRefreshToken,
	})

}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
