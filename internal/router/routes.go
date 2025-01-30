package router

import (
	"get_pet/internal/handler"
	"get_pet/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func BootstrapRouter(app *fiber.App, db *gorm.DB) {
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	api := app.Group("/api")

	userHandler := handler.NewUserHandler(db)
	authHandler := handler.NewAuthHandler(db)
	petHandler := handler.NewPetHandler(db)

	api.Post("/register", userHandler.Register)
	api.Post("/login", authHandler.Login)

	userRouter := api.Group("/user", middleware.AuthMiddleware)
	userRouter.Get("profile", userHandler.GetProfile)

	petRouter := api.Group("/pet", middleware.AuthMiddleware)
	petRouter.Post("/", petHandler.Create)
	petRouter.Get("/", petHandler.GetAll)
}
