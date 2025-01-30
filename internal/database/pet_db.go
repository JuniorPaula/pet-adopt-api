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

func (p *PetDB) GetByID(ID, userID int) (*model.Pet, error) {
	var pet model.Pet
	err := p.DB.Preload("Owner").Where("id = ? AND user_id = ?", ID, userID).First(&pet).Error
	if err != nil {
		return nil, err
	}

	return &pet, nil
}

func (p *PetDB) Update(pet *model.Pet) error {
	return p.DB.Where("id = ?", pet.ID).Updates(pet).Error
}
