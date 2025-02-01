package database

import (
	"get_pet/internal/model"

	"gorm.io/gorm"
)

type AdoptDB struct {
	DB *gorm.DB
}

func NewAdopt(db *gorm.DB) *AdoptDB {
	return &AdoptDB{
		DB: db,
	}
}

func (u *AdoptDB) Create(adopt *model.Adoption) error {
	return u.DB.Create(adopt).Error
}

func (r *AdoptDB) GetAdoptionsByUserID(userID uint) ([]model.Adoption, error) {
	var adoptions []model.Adoption
	err := r.DB.Preload("OldOwner").Preload("Adopter").Where("adopter_id = ?", userID).Preload("Pet").Find(&adoptions).Error
	if err != nil {
		return nil, err
	}
	return adoptions, nil
}
