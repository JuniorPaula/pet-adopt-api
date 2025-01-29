package main

import (
	"get_pet/internal/config"
	"get_pet/internal/router"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var PORT = "6969"

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error load .env file")
	}
}

func main() {
	db, err := config.ConnectDatabase()
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	router.BootstrapRouter(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = PORT
	}

	log.Fatal(app.Listen(":" + port))
}
