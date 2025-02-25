package handler

import (
	"get_pet/internal/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type VisitHandler struct {
	VisitDB database.VisitInterface
}

func NewVisitHandler(db *gorm.DB) *VisitHandler {
	return &VisitHandler{VisitDB: database.NewVisit(db)}
}

func (h *VisitHandler) GetUserVisits(c *fiber.Ctx) error {
	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	visits, err := h.VisitDB.GetVisitsByUserID(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	return c.Status(fiber.StatusOK).JSON(Response{Error: false, Message: "message", Data: visits})
}
