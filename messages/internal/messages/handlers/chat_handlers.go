package handlers

import (
	"messages/internal/messages/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CreateChat
func (h *MessagesHandlers) CreateChat(c echo.Context) error {
	var req domain.CreateChatRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	chat, err := h.service.CreateChat(c.Request().Context(), req.IsGroup)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, chat)
}

// DeleteChat
func (h *MessagesHandlers) DeleteChat(c echo.Context) error {
	var req domain.DeleteChatRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.service.DeleteChat(c.Request().Context(), req.ChatID); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "chat was delete")
}

// CreateGroupChat
func (h *MessagesHandlers) CreateGroupChat(c echo.Context) error {
	var req domain.CreateGroupChat

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	groupChat, err := h.service.CreateGroupChat(c.Request().Context(), &req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, groupChat)
}
