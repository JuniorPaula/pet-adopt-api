package database

import (
	"get_pet/internal/model"

	"gorm.io/gorm"
)

type VisitDB struct {
	DB *gorm.DB
}

func NewVisit(db *gorm.DB) *VisitDB {
	return &VisitDB{
		DB: db,
	}
}

func (u *VisitDB) Create(visit *model.Visit) error {
	return u.DB.Create(visit).Error
}

func (u *VisitDB) GetByPetID(petID int) (*model.Visit, error) {
	var v *model.Visit
	err := u.DB.Where("pet_id = ?", petID).First(&v).Error
	if err != nil {
		return nil, err
	}

	return v, nil
}
