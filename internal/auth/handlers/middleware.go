package handlers

import (
	"fmt"
	"messenger-app/internal/auth/store/generated"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	authHeader = "Authorization"
)

func (h *AuthHandler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get(authHeader)
		fmt.Println("Authorization header:", header)

		if header == "" {
			return c.JSON(http.StatusUnauthorized, map[string]any{"error": "missing header"})
		}

		if len(header) < 7 || header[:7] != "Bearer " {
			return c.JSON(http.StatusUnauthorized, map[string]any{"error": "wrong token format"})
		}

		token := header[7:]
		fmt.Println("Extracted token:", token)

		userID, err := h.service.ParseToken(c.Request().Context(), token, "access")

		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error": err.Error(),
			})
		}

		user, err := h.service.GetUserByID(c.Request().Context(), userID)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"error": "user not found",
			})
		}

		c.Set("user", user)
		c.Set("userId", userID)

		return next(c)
	}
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	userRaw := c.Get("user")
	dbUser, ok := userRaw.(generated.User)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid user data"})
	}

	return c.JSON(http.StatusOK, dbUser)
}
