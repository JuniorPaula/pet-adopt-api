package handler

import (
	"get_pet/internal/database"
	"get_pet/internal/dto"
	"get_pet/internal/model"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserDB        database.UserInterface
	UserDetailsDB database.UserDetailsInterface
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{UserDB: database.NewUser(db), UserDetailsDB: database.NewUserDetails(db)}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var body dto.CreateUserDto

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{Error: true, Message: ERRInternalServerError})
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
		if strings.Contains(err.Error(), ERRUniqueConstraint) {
			return c.Status(fiber.StatusConflict).JSON(Response{Error: true, Message: "user already exists"})
		}

		return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{Error: true, Message: "error to register user"})
	}

	return c.Status(fiber.StatusCreated).JSON(Response{
		Error:   false,
		Message: "user register on success",
		Data:    u,
	})
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	u, err := h.UserDB.GetByID(userID)
	if err != nil {
		if strings.Contains(err.Error(), ERRRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(Response{
				Error:   true,
				Message: "User not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Error:   true,
			Message: ERRInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{Error: false, Message: "success", Data: u})
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	var body dto.UpdateUserDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{Error: true, Message: ERRBadRequest})
	}

	u, err := h.UserDB.GetByID(userID)
	if err != nil {
		if strings.Contains(err.Error(), ERRRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(Response{
				Error:   true,
				Message: "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Error:   true,
			Message: ERRInternalServerError,
		})
	}

	updatedFileds := map[string]any{}

	if body.FirstName != "" {
		updatedFileds["FirstName"] = body.FirstName
	}
	if body.LastName != "" {
		updatedFileds["LastName"] = body.LastName
	}
	if body.Email != "" {
		updatedFileds["Email"] = body.Email
	}
	if body.Password != "" {
		updatedFileds["Password"], err = u.GenerateHashedPassword(body.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(Response{
				Error:   true,
				Message: ERRInternalServerError,
			})
		}
	}

	if body.Details.Address != "" {
		u.Details.Address = body.Details.Address
	}
	if body.Details.City != "" {
		u.Details.City = body.Details.City
	}
	if body.Details.ZipCode != "" {
		u.Details.ZipCode = body.Details.ZipCode
	}
	if body.Details.Phone != "" {
		u.Details.Phone = body.Details.Phone
	}
	if body.Details.Province != "" {
		u.Details.Province = body.Details.Province
	}
	u.Details.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	err = h.UserDB.Update(u, updatedFileds)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Error:   true,
			Message: ERRInternalServerError,
		})
	}

	err = h.UserDetailsDB.Update(&u.Details)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Error:   true,
			Message: ERRInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Error:   false,
		Message: "user updated on success",
		Data:    u,
	})
}

// DeleteProfile is a handler to disable a user account
func (h *UserHandler) DeleteProfile(c *fiber.Ctx) error {
	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	err = h.UserDB.SoftRemove(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Error:   true,
			Message: ERRInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Error:   false,
		Message: "Account disabled on success",
	})
}
