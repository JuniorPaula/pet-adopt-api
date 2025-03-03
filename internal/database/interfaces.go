package database

import "get_pet/internal/model"

// UserInterface is a model for user struct
type UserInterface interface {
	Create(user *model.User) error
	GetByEmail(string) (*model.User, error)
	GetByID(int) (*model.User, error)
}

// PetInterface is a model for pet struct
type PetInterface interface {
	Create(pet *model.Pet) error
	GetAll(page, limit int, sort string) ([]model.Pet, error)
	GetAllByUserID(userID, page, limit int, sort string) ([]model.Pet, error)
	GetByID(ID, userID int) (*model.Pet, error)
	Update(pet *model.Pet, newPet any) error
	UpdateImages(ID int, images []string) error
	UpdateAvailability(petID int, available bool) error
}

// VisitInterface is a model for visit struct
type VisitInterface interface {
	Create(visit *model.Visit) error
	GetVisitByID(ID int) (*model.Visit, error)
	FindVisitShceduledByAdopterID(petID int, userID uint) (*model.Visit, error)
	FindVisitShceduledByOnwerID(petID int, ownerID uint) (*model.Visit, error)
	GetVisitsByAdoperID(adopterID uint) ([]model.Visit, error)
	GetVisitsByOwnerID(ownerID uint) ([]model.Visit, error)
	Update(visit *model.Visit, newData any) error
	UpdateStatus(ID int, status string) error
}

// AdoptInterface is a model for adopt pet struct
type AdoptInterface interface {
	Create(adopt *model.Adoption) error
	GetAdoptionsByUserID(userID uint) ([]model.Adoption, error)
	FindAdoptionByPetIDAndAdopterID(petID int, adoptID uint) (*model.Adoption, error)
}
