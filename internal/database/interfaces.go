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
	Update(pet *model.Pet, newPet interface{}) error
	UpdateImages(ID int, images []string) error
	UpdateAvailability(petID int, available bool) error
}

// VisitInterface is a model for visit struct
type VisitInterface interface {
	Create(visit *model.Visit) error
	GetByPetID(petID int) (*model.Visit, error)
	Update(visit *model.Visit, newData interface{}) error
	UpdateStatus(ID int, status string) error
}

// AdoptInterface is a model for adopt pet struct
type AdoptInterface interface {
	Create(adopt *model.Adoption) error
}
