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
	api.Get("/pets", petHandler.GetAll)

	userRouter := api.Group("/users", middleware.AuthMiddleware)
	userRouter.Get("profile", userHandler.GetProfile)

	petRouter := api.Group("/pets", middleware.AuthMiddleware)
	petRouter.Post("/", petHandler.Create)
	petRouter.Get("/me", petHandler.GetAllByUserID)
	petRouter.Get("/:id", petHandler.GetByID)
	petRouter.Put("/:id", petHandler.Update)
	petRouter.Put("/:id/images", petHandler.UpdatePetImages)
	petRouter.Delete("/:id/images/:imageHash", petHandler.RemovePetImages)
}
