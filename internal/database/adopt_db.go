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
	err := r.DB.Preload("OldOwner.Details").Preload("Adopter.Details").Where("adopter_id = ?", userID).Preload("Pet").Find(&adoptions).Error
	if err != nil {
		return nil, err
	}
	return adoptions, nil
}

func (r *AdoptDB) FindAdoptionByPetIDAndAdopterID(petID int, adoptID uint) (*model.Adoption, error) {
	var adopt *model.Adoption
	err := r.DB.Preload("OldOwner.Details").Preload("Adopter.Details").Where("pet_id = ? AND adopter_id = ?", petID, adoptID).Preload("Pet").First(&adopt).Error
	if err != nil {
		return nil, err
	}
	return adopt, nil
}

func (r *AdoptDB) CountAdoptionsByOwnerID(ownerID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Adoption{}).Where("old_owner_id = ?", ownerID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *AdoptDB) GetAdoptionsByOldOwnerID(oldOwnerID uint) ([]model.Adoption, error) {
	var adoptions []model.Adoption
	err := r.DB.Preload("Adopter.Details").Limit(10).Order("id desc").Where("old_owner_id = ?", oldOwnerID).Preload("Pet").Find(&adoptions).Error
	if err != nil {
		return nil, err
	}
	return adoptions, nil
}
