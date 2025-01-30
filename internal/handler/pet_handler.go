package handler

import (
	"fmt"
	"get_pet/internal/database"
	"get_pet/internal/dto"
	"get_pet/internal/model"
	"get_pet/internal/util"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	var imagesPath []string
	var savedFiles []string

	form, err := c.MultipartForm()
	if err == nil {
		files := form.File["images"]
		for _, file := range files {
			ext := filepath.Ext(file.Filename)
			newFilename := fmt.Sprintf("%d-%s%s", time.Now().Unix(), util.GenerateRandomHash(12), ext)
			filePath := filepath.Join(uploadDir, newFilename)

			if err := c.SaveFile(file, filePath); err != nil {
				// Remove if not save
				if len(savedFiles) > 0 {
					for _, savedFile := range savedFiles {
						os.Remove(savedFile)
					}
				}

				return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: "Could not save file"})
			}

			imagesPath = append(imagesPath, "/"+filePath)
			savedFiles = append(savedFiles, filePath)
		}
	}

	pet := model.NewPet(userID, body.Age, body.Weight, body.Name, body.Size, body.Color, imagesPath)
	err = pet.ValidateFields()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{Error: true, Message: err.Error()})
	}

	err = h.PetDB.Create(pet)
	if err != nil {
		if len(savedFiles) > 0 {
			for _, savedFile := range savedFiles {
				os.Remove(savedFile)
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	return c.Status(fiber.StatusCreated).JSON(Response{
		Error:   false,
		Message: "pet created on success",
		Data:    pet,
	})
}

func (h *PetHandler) GetAll(c *fiber.Ctx) error {
	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	pets, err := h.PetDB.GetAll(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	return c.Status(fiber.StatusCreated).JSON(Response{
		Error:   false,
		Message: "success",
		Data:    pets,
	})
}

func (h *PetHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Error:   true,
			Message: ERRInternalServerError,
		})
	}

	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	pet, err := h.PetDB.GetByID(id, userID)
	if err != nil {
		if strings.Contains(err.Error(), ERRRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(Response{
				Error:   true,
				Message: "pet not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Error:   true,
			Message: ERRInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Error:   false,
		Message: "success",
		Data:    pet,
	})
}
