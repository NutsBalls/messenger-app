package handlers

import (
	"fmt"
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

	params := domain.CreateMessage(req)

	msg, err := h.service.CreateMessage(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp := domain.CreateMessageResponse{
		ChatID:    msg.ChatID,
		SenderID:  msg.SenderID,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
	}

	return c.JSON(http.StatusOK, resp)
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

	//TODO: асинхронно считывать и записывать
	resp := make([]domain.CreateMessageResponse, 0, len(msgs))
	for i, v := range msgs {
		resp[i] = domain.CreateMessageResponse{
			ChatID:    v.ChatID,
			SenderID:  v.SenderID,
			Content:   v.Content,
			CreatedAt: v.CreatedAt,
		}
	}

	return c.JSON(http.StatusOK, resp)
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

	return c.JSON(http.StatusOK, fmt.Sprintf("message %s was edited", req.MessageID.String()))
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

	return c.JSON(http.StatusOK, fmt.Sprintf("message %s was deleted", req.MessageID.String()))
}
