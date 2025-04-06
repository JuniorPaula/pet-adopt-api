package handler

import (
	"get_pet/internal/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AdoptionHandler struct {
	AdoptDB database.AdoptInterface
	VisitDB database.VisitInterface
}

func NewAdoptionHandler(db *gorm.DB) *AdoptionHandler {
	return &AdoptionHandler{AdoptDB: database.NewAdopt(db), VisitDB: database.NewVisit(db)}
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

func (h *AdoptionHandler) GetTotalAdoptionsAndVisitsByOwnerID(c *fiber.Ctx) error {
	ownerID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	adoptionCount, err := h.AdoptDB.CountAdoptionsByOwnerID(uint(ownerID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Error:   true,
			Message: "failed to count adoptions",
		})
	}

	visitCount, err := h.VisitDB.CountVisitsByOwnerID(uint(ownerID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Error:   true,
			Message: "failed to count visits",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Error:   false,
		Message: "message",
		Data: fiber.Map{
			"adoption_count":        adoptionCount,
			"visit_count":           visitCount,
			"visit_scheduled_count": visitCount,
		},
	})
}
