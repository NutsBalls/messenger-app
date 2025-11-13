package handlers

import (
	"fmt"
	"messages/internal/messages/domain/dto"
	"messages/internal/messages/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MessagesHandlers struct {
	service service.UseCase
}

func NewMessagesHandlers(service service.UseCase) *MessagesHandlers {
	return &MessagesHandlers{
		service: service,
	}
}

// CreateMessage
func (h *MessagesHandlers) CreateMessage(c echo.Context) error {
	var req dto.CreateMessageRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	msg, err := h.service.CreateMessage(c.Request().Context(), req.ToDomain())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, dto.ToMessageResponse(msg))
}

// GetMessages
func (h *MessagesHandlers) GetMessages(c echo.Context) error {
	var req dto.GetMessagesRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	msgs, err := h.service.GetMessages(c.Request().Context(), req.ChatID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//TODO: асинхронно считывать и записывать
	resps := make([]dto.MessageResponse, 0, len(msgs))
	for _, v := range msgs {
		resps = append(resps, dto.ToMessageResponse(v))
	}

	return c.JSON(http.StatusOK, resps)
}

// EditMessage
func (h *MessagesHandlers) EditMessage(c echo.Context) error {
	var req dto.EditMessageRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.service.EditMessage(c.Request().Context(), req.MessageID, req.NewContent); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("message %s was edited", req.MessageID.String()))
}

// DeleteMessage
func (h *MessagesHandlers) DeleteMessage(c echo.Context) error {
	var req dto.DeleteMessageRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.service.DeleteMessage(c.Request().Context(), req.MessageID); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("message %s was deleted", req.MessageID.String()))
}
