package handler

import (
	"get_pet/internal/database"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type VisitHandler struct {
	VisitDB database.VisitInterface
}

func NewVisitHandler(db *gorm.DB) *VisitHandler {
	return &VisitHandler{VisitDB: database.NewVisit(db)}
}

// GetAdopterVisits get all visits by adopter id
// user is able to see all visits that they have
// return a list of visits that the adopter has made to the pets
func (h *VisitHandler) GetAdopterVisits(c *fiber.Ctx) error {
	adopterID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	visits, err := h.VisitDB.GetVisitsByAdoperID(uint(adopterID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	return c.Status(fiber.StatusOK).JSON(Response{Error: false, Message: "message", Data: visits})
}

// GetOwnerVisits get all visits by owner id
// owner is able to see all visits that have been made to their pets
func (h *VisitHandler) GetOwnerVisits(c *fiber.Ctx) error {
	ownerID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	visits, err := h.VisitDB.GetVisitsByOwnerID(uint(ownerID))
	if err != nil {
		if strings.Contains(err.Error(), "(SQLSTATE 42703)") {
			return c.Status(fiber.StatusNotFound).JSON(Response{Error: true, Message: ERRRecordNotFound})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	return c.Status(fiber.StatusOK).JSON(Response{Error: false, Message: "message", Data: visits})
}

func (h *VisitHandler) UpdateVisitStatus(c *fiber.Ctx) error {
	visitID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: ERRBadRequest,
		})
	}

	ownerID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	visit, err := h.VisitDB.GetVisitByID(visitID)
	if err != nil {
		if strings.Contains(err.Error(), "(SQLSTATE 42703)") {
			return c.Status(fiber.StatusNotFound).JSON(Response{Error: true, Message: ERRRecordNotFound})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	// check if the owner is the owner of the pet
	if visit.OwnerPetId != uint(ownerID) {
		return c.Status(fiber.StatusUnauthorized).JSON(Response{Error: true, Message: "unauthorized"})
	}

	var data struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: ERRBadRequest,
		})
	}

	if data.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: "status is required",
		})
	}

	err = h.VisitDB.UpdateStatus(visitID, data.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	return c.Status(fiber.StatusOK).JSON(Response{Error: false, Message: "status updated"})
}
