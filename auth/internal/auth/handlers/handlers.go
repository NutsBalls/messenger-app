package handlers

import (
	"messenger-app/internal/auth/domain"
	"messenger-app/internal/auth/service"
	"messenger-app/pkg/hasher"
	"net/http"
	"strings"

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

func (h *AuthHandler) SignUp(c echo.Context) error {
	var req RegistredRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ReturnError(err, "bad request"))
	}

	if len([]byte(req.Email)) <= 7 && !strings.Contains(req.Email, "@") {
		return c.JSON(http.StatusBadRequest, ReturnMessage("wrong email for sign up"))
	}

	passHash, err := hasher.Hash(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ReturnError(err, "error in hash password"))
	}

	params := domain.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passHash,
	}

	user, err := h.service.SignUp(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ReturnError(err, "internal server error"))
	}

	resp := domain.CreateUserResponse{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req domain.LogInParams

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ReturnError(err, "bad request"))
	}

	if len([]byte(req.Email)) <= 7 && !strings.Contains(req.Email, "@") {
		return c.JSON(http.StatusBadRequest, ReturnMessage("wrong email"))
	}

	if req.Email == "" || req.PasswordHash == "" {
		return c.JSON(http.StatusBadRequest, ReturnMessage("wrong email/password"))
	}

	loginData, err := h.service.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ReturnError(err, "internal server error"))
	}

	return c.JSON(http.StatusOK, loginData)
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.Bind(&body); err != nil || body.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, ReturnError(err, "bad request"))
	}

	userID, err := h.service.ParseToken(c.Request().Context(), body.RefreshToken, "refresh")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ReturnError(err, "wrong refresh token"))
	}

	user, err := h.service.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, ReturnError(err, "user not found"))
	}

	if user.CryptedRefreshToken == nil {
		return c.JSON(http.StatusUnauthorized, ReturnError(err, "empty refresh token for user"))
	}

	hashedRefresh, err := hasher.HashToken(body.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ReturnError(err, "error while hashing token"))
	}

	if hashedRefresh != *user.CryptedRefreshToken {
		return c.JSON(http.StatusUnauthorized, ReturnError(err, "refresh tokem mismatch"))
	}

	newTokens, err := h.service.GenerateTokens(c.Request().Context(), user.ID.String())
	if err != nil {
		return c.JSON(http.StatusBadRequest, ReturnError(err, "internal server error"))
	}

	cryptedRefresh, err := hasher.HashToken(newTokens.Refresh)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ReturnError(err, "error while hashing token"))
	}

	if err := h.service.UpdateRefreshToken(c.Request().Context(), userID, &cryptedRefresh); err != nil {
		return c.JSON(http.StatusBadRequest, ReturnError(err, "update token error"))
	}

	return c.JSON(http.StatusOK, newTokens)

}
