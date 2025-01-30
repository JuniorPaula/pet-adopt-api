package database

import (
	"get_pet/internal/model"

	"gorm.io/gorm"
)

type PetDB struct {
	DB *gorm.DB
}

func NewPet(db *gorm.DB) *PetDB {
	return &PetDB{
		DB: db,
	}
}

func (u *PetDB) Create(pet *model.Pet) error {
	return u.DB.Create(pet).Error
}
