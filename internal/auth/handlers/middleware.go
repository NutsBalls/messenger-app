package handlers

import (
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

		if header == "" {
			return c.JSON(http.StatusUnauthorized, ReturnMessage("missing header"))
		}

		if len(header) < 7 || header[:7] != "Bearer " {
			return c.JSON(http.StatusUnauthorized, ReturnMessage("wrong header"))
		}

		token := header[7:]

		userID, err := h.service.ParseToken(c.Request().Context(), token, "access")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, ReturnError(err, "wrong token parse"))
		}

		user, err := h.service.GetUserByID(c.Request().Context(), userID)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, ReturnError(err, "user not found"))
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
		return c.JSON(http.StatusInternalServerError, ReturnMessage("invalid user data"))
	}

	return c.JSON(http.StatusOK, dbUser)
}
