package handlers

import (
	"messenger-app/internal/auth/service"

	"github.com/labstack/echo/v4"
)

func Router(e *echo.Echo, h *AuthHandler, uc service.AuthUsecase) {
	e.POST("/signup", h.SignUp)
	e.POST("/login", h.Login)

	api := e.Group("/api")
	api.Use(h.AuthMiddleware(uc))
	api.GET("/profile", h.GetProfile)
	api.POST("/auth/refresh", h.Refresh)
}
