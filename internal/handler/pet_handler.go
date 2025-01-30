package handler

import (
	"get_pet/internal/database"
	"get_pet/internal/dto"
	"get_pet/internal/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PetHandler struct {
	PetDB database.PetInterface
}

func NewPetHandler(db *gorm.DB) *PetHandler {
	return &PetHandler{PetDB: database.NewPet(db)}
}

func (h *PetHandler) Create(c *fiber.Ctx) error {
	var body dto.CreatePetDto

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	pet := model.NewPet(userID, body.Age, body.Weight, body.Name, body.Size, body.Color, body.Images)
	err = pet.ValidateFields()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{Error: true, Message: err.Error()})
	}

	err = h.PetDB.Create(pet)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	return c.Status(fiber.StatusCreated).JSON(Response{
		Error:   false,
		Message: "pet created on success",
		Data:    pet,
	})
}
