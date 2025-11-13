package handlers

import (
	"messenger-app/internal/auth/domain"
	"messenger-app/internal/auth/service"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *AuthHandler) AuthMiddleware(uc service.AuthUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, ReturnMessage("missing or invalid token"))
			}

			token := strings.TrimPrefix(header, "Bearer ")
			userID, err := uc.ParseToken(c.Request().Context(), token, "access")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, ReturnError(err, "wrong token"))
			}

			user, err := uc.GetUserByID(c.Request().Context(), userID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, ReturnError(err, "user not found"))
			}

			c.Set("user", user)
			return next(c)
		}
	}
}

func UserFromContext(c echo.Context) (domain.User, bool) {
	u, ok := c.Get("user").(domain.User)
	return u, ok
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	user, ok := UserFromContext(c)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ReturnMessage("invalid user data"))
	}

	return c.JSON(http.StatusOK, user)
}
