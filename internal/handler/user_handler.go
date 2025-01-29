package handler

import (
	"get_pet/internal/database"
	"get_pet/internal/dto"
	"get_pet/internal/model"

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
	var body dto.CreateUserDto

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{Error: true, Message: "internal server error"})
	}

	if body.Password != body.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(Response{Error: true, Message: "passwords doest match"})
	}

	u, err := model.NewUser(body.FirstName, body.LastName, body.Email, body.Password, false)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{Error: true, Message: "error to register user"})
	}

	err = h.UserDB.Create(u)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{Error: true, Message: "error to register user"})
	}

	return c.Status(fiber.StatusCreated).JSON(Response{
		Error:   false,
		Message: "user register on success",
		Data:    u,
	})
}
