package handler

import (
	"get_pet/internal/database"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthHandler struct {
	UserDB database.UserInterface
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{UserDB: database.NewUser(db)}
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body credentials

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	user, err := h.UserDB.GetByEmail(body.Email)
	if err != nil {
		if strings.Contains(err.Error(), ERRRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(Response{Error: true, Message: "unauthorized"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	if user.IsAccountActivated() {
		return c.Status(fiber.StatusUnauthorized).JSON(Response{Error: true, Message: "unauthorized"})
	}

	if !user.ValidatePassword(body.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(Response{
			Error:   true,
			Message: "unauthorized",
		})
	}

	claims := jwt.MapClaims{
		"sub":        user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"is_admin":   user.IsAdmin,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Error:   false,
		Message: "Login successful",
		Data:    map[string]interface{}{"user": user, "token": t},
	})
}
