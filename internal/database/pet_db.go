package database

import (
	"encoding/json"
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

func (p *PetDB) GetAll(page, limit int, sort string) ([]model.Pet, error) {
	var pets []model.Pet
	var err error

	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if page != 0 && limit != 0 {
		err = p.DB.Preload("Owner").Limit(limit).Offset((page - 1) * limit).Order("created_at " + sort).Find(&pets).Error
	} else {
		err = p.DB.Preload("Owner").Find(&pets).Error
	}

	return pets, err
}

func (p *PetDB) GetAllByUserID(userID, page, limit int, sort string) ([]model.Pet, error) {
	var pets []model.Pet
	var err error

	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if page != 0 && limit != 0 {
		err = p.DB.Preload("Owner").Where("user_id = ?", userID).Limit(limit).Offset((page - 1) * limit).Order("created_at " + sort).Find(&pets).Error
	} else {
		err = p.DB.Preload("Owner").Where("user_id = ?", userID).Find(&pets).Error
	}

	return pets, err
}

func (p *PetDB) GetByID(ID, userID int) (*model.Pet, error) {
	var pet model.Pet
	err := p.DB.Preload("Owner").Where("id = ? AND user_id = ?", ID, userID).First(&pet).Error
	if err != nil {
		return nil, err
	}

	return &pet, nil
}

func (p *PetDB) Update(pet *model.Pet, newPet interface{}) error {
	return p.DB.Model(&pet).Where("id = ?", pet.ID).Updates(newPet).Error
}

func (p *PetDB) UpdateImages(ID int, images []string) error {
	imagesJSON, err := json.Marshal(images)
	if err != nil {
		return err
	}

	return p.DB.Model(&model.Pet{}).Where("id = ?", ID).Update("images", imagesJSON).Error
}

func (p *PetDB) UpdateAvailability(petID int, available bool) error {
	return p.DB.Model(&model.Pet{}).Where("id = ?", petID).Update("available", available).Error
}
