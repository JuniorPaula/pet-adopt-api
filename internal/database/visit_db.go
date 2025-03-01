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

func (vd *VisitDB) Create(visit *model.Visit) error {
	return vd.DB.Create(visit).Error
}

// GetVisitByPetIDAndUserID get visit by pet id and user id
// return a visit to the user
func (vd *VisitDB) GetVisitByPetIDAndUserID(petID int, userID uint) (*model.Visit, error) {
	var v *model.Visit
	err := vd.DB.Where("pet_id = ? AND user_id = ?", petID, userID).First(&v).Error
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (vd *VisitDB) Update(visit *model.Visit, newData interface{}) error {
	return vd.DB.Model(&visit).Where("id = ?", visit.ID).Updates(newData).Error
}

func (vd *VisitDB) UpdateStatus(ID int, status string) error {
	return vd.DB.Model(&model.Visit{}).Where("id = ?", ID).Update("status", status).Error
}

// GetVisitsByUserID get all visits by user id
// user is able to see all visits that they have
// return a list of visits that the user has made to the pets
func (vd *VisitDB) GetVisitsByUserID(userID uint) ([]model.Visit, error) {
	var visits []model.Visit
	err := vd.DB.Where("user_id = ?", userID).Preload("Pet", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Owner")
	}).Find(&visits).Error

	if err != nil {
		return nil, err
	}

	return visits, nil
}

// GetVisitsByOwnerID get all visits by owner id
// owner is able to see all visits that have been made to their pets
// return a list of visits that have been made to the owner's pets
func (vd *VisitDB) GetVisitsByOwnerID(ownerID uint) ([]model.Visit, error) {
	var visits []model.Visit
	err := vd.DB.Where("owner_pet_id = ? and status = 'pending'", ownerID).Preload("User").Preload("Pet").Find(&visits).Error

	if err != nil {
		return nil, err
	}

	return visits, nil
}
