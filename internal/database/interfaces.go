package database

import "get_pet/internal/model"

type UserInterface interface {
	Create(user *model.User) error
	GetByEmail(string) (*model.User, error)
	GetByID(int) (*model.User, error)
}

type PetInterface interface {
	Create(pet *model.Pet) error
}
