package handler

import (
	"get_pet/internal/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AdoptionHandler struct {
	AdoptDB database.AdoptInterface
}

func NewAdoptionHandler(db *gorm.DB) *AdoptionHandler {
	return &AdoptionHandler{AdoptDB: database.NewAdopt(db)}
}

// GetUserAdoptions get all adoptions by user id (adopter)
func (h *AdoptionHandler) GetUserAdoptions(c *fiber.Ctx) error {
	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	adoptions, err := h.AdoptDB.GetAdoptionsByUserID(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	return c.Status(fiber.StatusOK).JSON(Response{Error: false, Message: "message", Data: adoptions})
}

// GetOneAdoption get one adoption by user id (adopter)
func (h *AdoptionHandler) GetOneAdoption(c *fiber.Ctx) error {
	petID, err := strconv.Atoi(c.Params("petID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{Error: true, Message: "invalid pet id"})
	}

	adopterID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	adoptions, err := h.AdoptDB.FindAdoptionByPetIDAndAdopterID(petID, uint(adopterID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	return c.Status(fiber.StatusOK).JSON(Response{Error: false, Message: "message", Data: adoptions})
}
