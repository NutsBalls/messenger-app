package main

import (
	"context"
	"errors"
	"log"
	"messages/internal/messages/handlers"
	"messages/internal/messages/service"
	"messages/internal/messages/store"
	"messages/pkg/config"
	"messages/pkg/database"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {

	cfg := config.Load()
	db := database.NewConn(cfg.DBURL)

	repo := store.NewMessagesRepository(db)
	svc := service.NewMessagesService(repo)

	defer db.Close()
	e := echo.New()

	h := handlers.NewHandlers(handlers.NewMessagesHandlers(svc), handlers.NewChatsHandlers(svc), handlers.NewMembersHandlers(svc))

	api := e.Group("/api")
	h.RegisterRoutes(api)

	port := ":" + cfg.Port
	go func() {
		if err := e.Start(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("shutting down server: %v", err)
		}
	}()
	log.Printf("server started on %s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
}
