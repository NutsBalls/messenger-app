package handlers

import (
	"fmt"
	"messages/internal/messages/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

// AddUserToChat
func (h *MessagesHandlers) AddUserToChat(c echo.Context) error {
	var req domain.UserChatRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	params := domain.UserChat(req)

	if err := h.service.AddUserToChat(c.Request().Context(), params); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("user %s added in chat %s", req.UserID.String(), req.ChatID.String()))

}

// RemoveUserFromChat
func (h *MessagesHandlers) RemoveUserFromChat(c echo.Context) error {
	var req domain.UserChatRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	params := domain.UserChat(req)

	if err := h.service.RemoveUserFromChat(c.Request().Context(), params); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("user %s was removed from chat %s", req.UserID.String(), req.ChatID.String()))
}

// GetChatMembers
func (h *MessagesHandlers) GetChatMembers(c echo.Context) error {
	var req domain.UserChatRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	members, err := h.service.GetChatMembers(c.Request().Context(), req.ChatID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, members)
}

// GetUserChats
func (h *MessagesHandlers) GetUserChats(c echo.Context) error {
	var req domain.UserChatRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	chats, err := h.service.GetUserChats(c.Request().Context(), req.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//TODO: асинхронно считвывать и записывать
	resp := make([]domain.ChatResponse, 0, len(chats))
	for i, v := range chats {
		resp[i] = domain.ChatResponse{
			Name:    v.Name,
			IsGroup: v.IsGroup,
		}
	}

	return c.JSON(http.StatusOK, resp)
}
