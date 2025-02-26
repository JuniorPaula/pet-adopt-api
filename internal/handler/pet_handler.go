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
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type PetHandler struct {
	PetDB   database.PetInterface
	VisitDB database.VisitInterface
	AdoptDB database.AdoptInterface
}

func NewPetHandler(db *gorm.DB) *PetHandler {
	return &PetHandler{
		PetDB:   database.NewPet(db),
		VisitDB: database.NewVisit(db),
		AdoptDB: database.NewAdopt(db),
	}
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
	var (
		page  = c.Query("page")
		limit = c.Query("limit")
		sort  = c.Query("sort")
	)

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}

	pets, err := h.PetDB.GetAll(pageInt, limitInt, sort)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	return c.Status(fiber.StatusCreated).JSON(Response{
		Error:   false,
		Message: "success",
		Data:    pets,
	})
}

func (h *PetHandler) GetAllByUserID(c *fiber.Ctx) error {
	var (
		page  = c.Query("page")
		limit = c.Query("limit")
		sort  = c.Query("sort")
	)

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}

	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	pets, err := h.PetDB.GetAllByUserID(userID, pageInt, limitInt, sort)
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
	petId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Error:   true,
			Message: ERRInternalServerError,
		})
	}

	pet, err := h.PetDB.GetByID(petId, 0)
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

func (h *PetHandler) GetMyPetByID(c *fiber.Ctx) error {
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

func (h *PetHandler) Update(c *fiber.Ctx) error {
	petId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{Error: true, Message: "invalid pet id"})
	}

	var body dto.UpdatePetDto

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

	pet, err := h.PetDB.GetByID(petId, userID)
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

	updateFields := map[string]interface{}{}

	if body.Name != "" {
		updateFields["name"] = body.Name
	}
	if body.Age > 0 {
		updateFields["Age"] = body.Age
	}
	if body.Weight > 0 {
		updateFields["Weight"] = body.Weight
	}
	if body.Color != "" {
		updateFields["Color"] = body.Color
	}
	if body.Size != "" {
		updateFields["Size"] = body.Size
	}
	if body.Available != nil {
		updateFields["Available"] = *body.Available
	}

	err = h.PetDB.Update(pet, updateFields)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: ERRInternalServerError})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Error:   false,
		Message: "pet updated on success",
		Data:    pet,
	})
}

func (h *PetHandler) UpdatePetImages(c *fiber.Ctx) error {
	petId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: "invalid pet id",
		})
	}

	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	pet, err := h.PetDB.GetByID(petId, userID)
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

	pet.Images = append(pet.Images, imagesPath...)
	err = h.PetDB.UpdateImages(petId, pet.Images)
	if err != nil {
		for _, savedFile := range savedFiles {
			os.Remove(savedFile)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: "error on upload images"})
	}

	return c.Status(fiber.StatusOK).JSON(Response{Error: false, Message: "upload image succeed"})
}

func (h *PetHandler) RemovePetImages(c *fiber.Ctx) error {
	imageHash := c.Params("imageHash")
	petId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: "invalid pet id",
		})
	}

	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	pet, err := h.PetDB.GetByID(petId, userID)
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

	var newImages []string
	for _, img := range pet.Images {
		parts := strings.Split(img, "-")

		hash := strings.Split(parts[1], ".")[0]
		if imageHash != hash {
			newImages = append(newImages, img)
		}
	}

	err = h.PetDB.UpdateImages(petId, newImages)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{Error: true, Message: "error on upload images"})
	}

	return c.Status(fiber.StatusOK).JSON(Response{Error: false, Message: "upload image succeed"})
}

// ScheduleVisit is a handler that pet visit schedule flow,
// the user who wants to adopt the pet schedules a visit to the pet owner
func (h *PetHandler) ScheduleVisit(c *fiber.Ctx) error {
	petId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{Error: true, Message: "invalid pet id"})
	}

	var visitData struct {
		PetOwnerID int `json:"owner_id"`
	}

	if err := c.BodyParser(&visitData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{Error: true, Message: "Invalid request"})
	}

	// this id is the user that to be adopt
	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	visit, err := h.VisitDB.GetVisitByPetIDAndUserID(petId, uint(userID))
	if err != nil {
		if strings.Contains(err.Error(), ERRInternalServerError) {
			return c.Status(fiber.StatusInternalServerError).JSON(Response{
				Error:   true,
				Message: ERRInternalServerError,
			})
		}
	}

	if visit != nil && visit.UserID == uint(userID) {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: "visit already exists",
		})
	}

	pet, err := h.PetDB.GetByID(petId, visitData.PetOwnerID)
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

	if !pet.Available {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: "pet is not available",
		})
	}

	if int(pet.UserID) == userID {
		return c.Status(fiber.StatusNotFound).JSON(Response{
			Error:   true,
			Message: "could not scheduler visit on your owner pet",
		})
	}

	newVisit := model.NewVisit(userID, petId, visitData.PetOwnerID, "pending")

	err = h.VisitDB.Create(newVisit)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(Response{
			Error:   true,
			Message: "error on scheduler visit pet",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(Response{
		Error:   false,
		Message: "visit scheduler on success",
	})
}

// GetVisitSchedule is a handler that get pet visit schedule flow,
// return the pet data if the visit schedule is available
func (h *PetHandler) GetVisitSchedule(c *fiber.Ctx) error {
	petId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{Error: true, Message: "invalid pet id"})
	}

	// this'id owner user pet
	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	// TODO: return pet data on visit schedule, make left join query
	visit, err := h.VisitDB.GetVisitByPetIDAndUserID(petId, uint(userID))
	if err != nil {
		if strings.Contains(err.Error(), ERRInternalServerError) {
			return c.Status(fiber.StatusInternalServerError).JSON(Response{
				Error:   true,
				Message: ERRInternalServerError,
			})
		}
	}

	if visit != nil && visit.UserID == uint(userID) {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: "visit already exists",
		})
	}

	// TODO: return pet data on visit schedule
	return c.Status(fiber.StatusOK).JSON(Response{
		Error:   false,
		Message: "success",
	})
}

// ConfirmAdopt is a handler that pet adoption flow,
// the user who owns the pet who completes the flow based on the visa table for the pet
func (h *PetHandler) ConfirmAdopt(c *fiber.Ctx) error {
	petId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{Error: true, Message: "invalid pet id"})
	}

	// this'id owner user pet
	userID, err := getUserIdFromCtx(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Error:   true,
			Message: err.Error(),
		})
	}

	pet, err := h.PetDB.GetByID(petId, userID)
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

	visit, err := h.VisitDB.GetVisitByPetIDAndUserID(pet.ID, uint(userID))
	if err != nil {
		if strings.Contains(err.Error(), ERRRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(Response{
				Error:   true,
				Message: fmt.Sprintf("`%s` doesn't have visit schedule!", pet.Name),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Error:   true,
			Message: ERRInternalServerError,
		})
	}

	// check if the user scheduling the visit is the same owner of the pet
	if visit.UserID == pet.UserID {
		return c.Status(fiber.StatusNotFound).JSON(Response{
			Error:   true,
			Message: "You can't adopt your owner pet!",
		})
	}

	petID, oldOwnerID, newOwnerID := pet.ID, pet.UserID, visit.UserID
	adoptData := model.NewAdoption(uint(petID), oldOwnerID, newOwnerID)

	err = h.AdoptDB.Create(adoptData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Error:   true,
			Message: ERRInternalServerError,
		})
	}

	// complete visit status
	if err := h.VisitDB.UpdateStatus(visit.ID, "completed"); err != nil {
		log.Errorf("error on update status visit: err %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Error:   true,
			Message: ERRInternalServerError,
		})
	}

	// update available pet field
	err = h.PetDB.UpdateAvailability(pet.ID, false)
	if err != nil {
		log.Errorf("error on update status pet: err %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Error:   true,
			Message: ERRInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Error:   false,
		Message: "Congratulations adoption completed successfully",
	})
}
