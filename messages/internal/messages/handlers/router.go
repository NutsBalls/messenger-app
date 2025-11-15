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
	messages := e.Group("/messages")
	messages.POST("/create", h.Messages.CreateMessage)
	messages.GET("/all", h.Messages.GetMessages)
	messages.PATCH("/edit", h.Messages.EditMessage)
	messages.DELETE("/delete", h.Messages.DeleteMessage)

	chats := e.Group("/chats")
	chats.POST("/create", h.Chats.CreateChat)
	chats.DELETE("/delete", h.Chats.DeleteChat)
	chats.POST("/crete_group", h.Chats.CreateGroupChat)

	members := e.Group("/members")
	members.POST("/add", h.Members.AddUserToChat)
	members.DELETE("/remove", h.Members.RemoveUserFromChat)
	members.GET("/users", h.Members.GetChatMembers)
	members.GET("/chats", h.Members.GetUserChats)
}
