package handlers

import (
	"messenger-app/internal/auth/domain"
	"messenger-app/internal/auth/service"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	authHeader = "Authorization"
)

type AuthUsecase interface {
	ParseToken(ctx echo.Context, token string, kind string) (uuid.UUID, error)
	GetUserByID(ctx echo.Context, id uuid.UUID) (domain.User, error)
}

func AuthMiddleware(uc service.AuthUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Request().Header.Get("Authorization")
			if h == "" || !strings.HasPrefix(h, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing or invalid auth header"})
			}
			token := strings.TrimPrefix(h, "Bearer ")

			userID, err := uc.ParseToken(c.Request().Context(), token, "access")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, ReturnError(err, "wrong token parse"))
			}

			user, err := uc.GetUserByID(c.Request().Context(), userID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, ReturnError(err, "user not found"))
			}

			c.Set("user", user)
			c.Set("userId", userID)

			return next(c)
		}
	}
}

func UserFromContext(c echo.Context) (domain.User, bool) {
	u, ok := c.Get("user").(domain.User)
	return u, ok
}
