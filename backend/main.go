package main

import (
	"project/config"
	"project/database"
	"project/server"

	"github.com/gofiber/fiber/v3/log"
)

func main() {
	db, err := database.New(config.C.DSN)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	app := server.New(db)
	if err := app.Listen(config.C.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
