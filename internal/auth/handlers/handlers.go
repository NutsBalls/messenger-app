package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func login(c echo.Context) error {
	var req User

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user := createUser(req.Name, req.Email, req.Password)

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Создан пользователь",
		"ID":      user.ID,
		"email":   user.Email,
	})
}
