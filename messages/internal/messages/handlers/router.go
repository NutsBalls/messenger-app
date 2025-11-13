package handlers

import "github.com/labstack/echo/v4"

type Handlers struct {
	Messages *MessagesHandlers
	Chats    *ChatsHandlers
	Members  *MembersHandlers
}

func NewHandlers(messages *MessagesHandlers, chats *ChatsHandlers, members *MembersHandlers) *Handlers {
	return &Handlers{
		Messages: messages,
		Chats:    chats,
		Members:  members,
	}
}

func (h *Handlers) RegisterRoutes(e *echo.Group) {
	e.POST("/create", h.Messages.CreateMessage)
	e.GET("/all", h.Messages.GetMessages)
	e.PATCH("/edit", h.Messages.EditMessage)
	e.DELETE("/delete", h.Messages.DeleteMessage)

	e.POST("/create", h.Chats.CreateChat)
	e.DELETE("/delete", h.Chats.DeleteChat)
	e.POST("/crete_group", h.Chats.CreateGroupChat)

	e.POST("/add", h.Members.AddUserToChat)
	e.DELETE("/remove", h.Members.RemoveUserFromChat)
	e.GET("/users", h.Members.GetChatMembers)
	e.GET("/chats", h.Members.GetUserChats)
}
