package main

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.POST("/signup", signUp)
	e.POST("/user", handlers.login)

	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}

func signUp(c echo.Context) error {
	return c.String(http.StatusOK, "added user")
}
