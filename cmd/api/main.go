package main

import (
	"get_pet/internal/config"
	"get_pet/internal/router"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	db, err := config.ConnectDatabase()
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	router.BootstrapRouter(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "6969"
	}

	log.Fatal(app.Listen(":" + port))
}
