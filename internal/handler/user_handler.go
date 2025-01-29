package handler

import (
	"fmt"
	"get_pet/internal/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{UserDB: database.NewUser(db)}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "internal server error"})
	}

	fmt.Println("payload: ", data)
	return nil
}
