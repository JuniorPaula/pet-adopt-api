package router

import (
	"get_pet/internal/handler"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BootstrapRouter(app *fiber.App, db *gorm.DB) {

	userHandler := handler.NewUserHandler(db)
	authHandler := handler.NewAuthHandler(db)

	app.Post("/register", userHandler.Register)
	app.Post("/login", authHandler.Login)
}
