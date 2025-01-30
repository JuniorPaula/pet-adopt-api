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

func (p *PetDB) Create(pet *model.Pet) error {
	return p.DB.Preload("Owner").Create(pet).Error
}

func (p *PetDB) GetAll(userID int) ([]model.Pet, error) {
	var pets []model.Pet
	err := p.DB.Preload("Owner").Where("user_id = ?", userID).Find(&pets).Error
	if err != nil {
		return pets, err
	}

	return pets, nil
}
