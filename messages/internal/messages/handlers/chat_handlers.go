package handlers

import (
	"messages/internal/messages/domain/dto"
	"messages/internal/messages/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ChatsHandlers struct {
	service service.UseCase
}

func NewChatsHandlers(s service.UseCase) *ChatsHandlers {
	return &ChatsHandlers{service: s}
}

// CreateChat
func (h *ChatsHandlers) CreateChat(c echo.Context) error {
	var req dto.CreateChatRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	chat, err := h.service.CreateChat(c.Request().Context(), req.IsGroup)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, dto.ToChatResponse(chat))
}

// DeleteChat
func (h *ChatsHandlers) DeleteChat(c echo.Context) error {
	var req dto.DeleteChatRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.service.DeleteChat(c.Request().Context(), req.ChatID); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "chat was delete")
}

// CreateGroupChat
func (h *ChatsHandlers) CreateGroupChat(c echo.Context) error {
	var req dto.CreateGroupChatRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	groupChat, err := h.service.CreateGroupChat(c.Request().Context(), &req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, dto.ToChatResponse(groupChat))
}
