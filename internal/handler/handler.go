package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var (
	ERRUniqueConstraint    = "SQLSTATE 23505"
	ERRInternalServerError = "Internal Server Error"
	ERRRecordNotFound      = "record not found"
	ERRBadRequest          = "bad request"
)

// Response to formalize http response data
type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func getUserIdFromCtx(c *fiber.Ctx) (int, error) {
	userCtx := c.Locals("user")
	if userCtx == "" {
		return 0, errors.New("anauthoried")
	}

	userMap, ok := userCtx.(fiber.Map)
	if !ok {
		return 0, errors.New("invalid user format")
	}
	idInterface := userMap["id"]

	var userID int
	switch v := idInterface.(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	case string:
		parsedID, err := strconv.Atoi(v)
		if err != nil {
			return 0, errors.New("invalid user ID format")
		}
		userID = parsedID
	default:
		return 0, errors.New("unexpected user ID type")
	}
	return userID, nil
}
