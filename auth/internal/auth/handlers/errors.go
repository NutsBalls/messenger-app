package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type MessageResponse struct {
	Message string
}

func ReturnError(err error, message string) ErrResponse {
	return ErrResponse{Error: err.Error(), Message: message}
}

func ReturnMessage(message string) MessageResponse {
	return MessageResponse{
		Message: message,
	}
}

func InternalError(c echo.Context, err error) error {
	log.Printf("internal error: %v", err)
	return c.JSON(http.StatusInternalServerError, ReturnMessage("internal error"))
}
