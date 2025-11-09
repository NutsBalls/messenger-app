package handlers

import (
	"messages/internal/messages/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CreateMessage
func (h *MessagesHandlers) CreateMessage(c echo.Context) error {
	var req domain.CreateMessageRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	msg, err := h.service.CreateMessage(c.Request().Context(), req.ChatID, req.SenderID, req.Content)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, msg)
}

// GetMessages
func (h *MessagesHandlers) GetMessages(c echo.Context) error {
	var req domain.GetMessagesRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	msgs, err := h.service.GetMessages(c.Request().Context(), req.ChatID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, msgs)
}

// EditMessage
func (h *MessagesHandlers) EditMessage(c echo.Context) error {
	var req domain.EditMessageRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.service.EditMessage(c.Request().Context(), req.MessageID, req.NewContent); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "message was edited")
}

// DeleteMessage
func (h *MessagesHandlers) DeleteMessage(c echo.Context) error {
	var req domain.DeleteMessageRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.service.DeleteMessage(c.Request().Context(), req.MessageID); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "message was deleted")
}
