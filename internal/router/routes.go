package router

import (
	"get_pet/internal/handler"
	"get_pet/internal/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func BootstrapRouter(app *fiber.App, db *gorm.DB) {
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	/* setup cors*/
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	/* setup statics images files*/
	app.Static("/uploads", "./uploads", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Index:         "",
		CacheDuration: 10 * time.Minute,
	})

	api := app.Group("/api")

	userHandler := handler.NewUserHandler(db)
	authHandler := handler.NewAuthHandler(db)
	petHandler := handler.NewPetHandler(db)
	adoptionHandler := handler.NewAdoptionHandler(db)

	/* publics routes */
	api.Post("/register", userHandler.Register)
	api.Post("/login", authHandler.Login)
	api.Get("/pets", petHandler.GetAll)

	/* users routes */
	userRouter := api.Group("/users", middleware.AuthMiddleware)
	userRouter.Get("profile", userHandler.GetProfile)

	/* pets routes */
	petRouter := api.Group("/pets", middleware.AuthMiddleware)
	petRouter.Post("/", petHandler.Create)
	petRouter.Get("/me", petHandler.GetAllByUserID)
	petRouter.Get("/:id", petHandler.GetByID)
	petRouter.Get("/:id/me", petHandler.GetMyPetByID)
	petRouter.Put("/:id", petHandler.Update)
	petRouter.Put("/:id/images", petHandler.UpdatePetImages)
	petRouter.Delete("/:id/images/:imageHash", petHandler.RemovePetImages)
	petRouter.Post("/:id/scheduler", petHandler.ScheduleVisit)
	petRouter.Get("/:id/scheduler", petHandler.GetVisitSchedule)
	petRouter.Post("/:id/adopt", petHandler.ConfirmAdopt)

	/* adoptions routes*/
	adoptRouter := api.Group("/adopts", middleware.AuthMiddleware)
	adoptRouter.Get("/", adoptionHandler.GetUserAdoptions)
}
