package main

import (
	"get_pet/internal/router"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func main() {
	app := fiber.New()

	router.BootstrapRouter(app, &gorm.DB{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "6969"
	}

	log.Fatal(app.Listen(":" + port))
}
