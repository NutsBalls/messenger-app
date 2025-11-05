package main

import (
	"log"
	"messenger-app/internal/auth/handlers"
	"messenger-app/internal/auth/service"
	"messenger-app/internal/auth/store"
	"messenger-app/pkg/config"
	"messenger-app/pkg/database"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.Load()
	db := database.NewConn(cfg.DBURL)

	store := store.NewAuthRepository(db)
	service := service.NewAuthService(store, cfg.JWTSecret)
	handler := handlers.NewAuthHandler(service)

	e := echo.New()

	e.POST("/signup", handler.SignUp)
	e.POST("/login", handler.Login)

	auth := e.Group("/auth")
	auth.Use(handler.AuthMiddleware)
	auth.GET("/profile", handler.GetProfile)
	auth.GET("/refresh", handler.Refresh)

	port := ":" + cfg.Port
	log.Printf("Starting server on port %s", port)
	if err := e.Start(port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
