package main

import (
	"log"
	"messages/pkg/config"
	"messages/pkg/database"

	"github.com/labstack/echo/v4"
)

func main() {

	cfg := config.Load()
	db := database.NewConn(cfg.DBURL)

	defer db.Close()
	e := echo.New()

	port := ":" + cfg.Port
	log.Printf("Messages service start on port %v", port)
	if err := e.Start(port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
