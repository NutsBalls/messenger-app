package main

import (
	"context"
	"errors"
	"log"
	"messenger-app/internal/auth/handlers"
	"messenger-app/internal/auth/service"
	"messenger-app/internal/auth/store"
	"messenger-app/pkg/config"
	"messenger-app/pkg/database"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	db := database.NewConn(cfg.DBURL)

	repo := store.NewAuthStore(db)
	svc := service.NewAuthService(repo, cfg.JWTSecret)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := handlers.NewAuthHandler(svc)
	handlers.Router(e, h, svc)

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
