package handlers

import (
	"fmt"
	"messages/internal/messages/domain/dto"
	"messages/internal/messages/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MembersHandlers struct {
	service service.UseCase
}

func NewMembersHandlers(s service.UseCase) *MembersHandlers {
	return &MembersHandlers{service: s}
}

// AddUserToChat
func (h *MembersHandlers) AddUserToChat(c echo.Context) error {
	var req dto.UserChatRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	if err := h.service.AddUserToChat(c.Request().Context(), req.ToDomain()); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("user %s added in chat %s", req.UserID.String(), req.ChatID.String()))

}

// RemoveUserFromChat
func (h *MembersHandlers) RemoveUserFromChat(c echo.Context) error {
	var req dto.UserChatRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	if err := h.service.RemoveUserFromChat(c.Request().Context(), req.ToDomain()); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("user %s was removed from chat %s", req.UserID.String(), req.ChatID.String()))
}

// GetChatMembers
func (h *MembersHandlers) GetChatMembers(c echo.Context) error {
	var req dto.UserChatRequest

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
func (h *MembersHandlers) GetUserChats(c echo.Context) error {
	var req dto.UserChatRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	chats, err := h.service.GetUserChats(c.Request().Context(), req.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//TODO: асинхронно считвывать и записывать
	resp := make([]dto.ChatResponse, 0, len(chats))
	for _, v := range chats {
		resp = append(resp, dto.ToChatResponse(v))
	}

	return c.JSON(http.StatusOK, resp)
}
