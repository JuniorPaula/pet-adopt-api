package router

import (
	"get_pet/internal/handler"
	"get_pet/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BootstrapRouter(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")

	userHandler := handler.NewUserHandler(db)
	authHandler := handler.NewAuthHandler(db)

	api.Post("/register", userHandler.Register)
	api.Post("/login", authHandler.Login)

	userRouter := api.Group("/user", middleware.AuthMiddleware)
	userRouter.Get("profile", func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": true, "message": "anauthoried"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "message": "user log in", "data": user})
	})
}
