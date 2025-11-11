package handlers

import (
	"messages/internal/messages/service"

	"github.com/labstack/echo/v4"
)

func Setup() {

	e := echo.New()

	h := MessagesHandlers{
		service: &service.MessagesService{},
	}
	// e.POST("/signup", handler.SignUp)
	// e.POST("/login", handler.Login)

	// auth := e.Group("/auth")
	// auth.Use(handler.AuthMiddleware)
	// auth.GET("/profile", handler.GetProfile)
	// auth.GET("/refresh", handler.Refresh)

	// port := ":" + cfg.Port
	// log.Printf("Starting server on port %s", port)
	// if err := e.Start(port); err != nil {
	// 	log.Fatal("Failed to start server: ", err)
	// }

	message := e.Group("/message")
	message.POST("/create", h.CreateMessage)
	message.GET("/all", h.GetMessages)
	message.PATCH("/edit", h.EditMessage)
	message.DELETE("/delete", h.DeleteMessage)

	chat := e.Group("/chat")
	chat.POST("/create", h.CreateChat)
	chat.DELETE("/delete", h.DeleteChat)
	chat.POST("/crete_group", h.CreateGroupChat)

	member := e.Group("/member")
	member.POST("/add", h.AddUserToChat)
	member.DELETE("/remove", h.RemoveUserFromChat)
	member.GET("/users", h.GetChatMembers)
	member.GET("/chats", h.GetUserChats)

	// AddUserToChat
	// RemoveUserFromChat
	// GetChatMembers
	// GetUserChats
}
