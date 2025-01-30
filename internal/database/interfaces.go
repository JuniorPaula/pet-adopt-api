package database

import "get_pet/internal/model"

type UserInterface interface {
	Create(user *model.User) error
	GetByEmail(string) (*model.User, error)
	GetByID(int) (*model.User, error)
}

type PetInterface interface {
	Create(pet *model.Pet) error
	GetAll(page, limit int, sort string) ([]model.Pet, error)
	GetAllByUserID(userID, page, limit int, sort string) ([]model.Pet, error)
	GetByID(ID, userID int) (*model.Pet, error)
	Update(pet *model.Pet) error
}
