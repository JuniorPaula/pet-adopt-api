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

// GetVisitByID get visit by id
// return a visit to the user
func (vd *VisitDB) GetVisitByID(ID int) (*model.Visit, error) {
	var v *model.Visit
	err := vd.DB.Where("id = ?", ID).First(&v).Error
	if err != nil {
		return nil, err
	}

	return v, nil
}

// FindVisitShceduledByAdopterID find the visit by pet id and adopter id
// return a visit to the user
func (vd *VisitDB) FindVisitShceduledByAdopterID(petID int, adopterID uint) (*model.Visit, error) {
	var v *model.Visit
	err := vd.DB.Where("pet_id = ? AND user_id = ?", petID, adopterID).First(&v).Error
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (vd *VisitDB) FindVisitShceduledByOnwerID(petID int, ownerID uint) (*model.Visit, error) {
	var v *model.Visit
	err := vd.DB.Where("pet_id = ? AND owner_pet_id = ?", petID, ownerID).First(&v).Error
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (vd *VisitDB) Update(visit *model.Visit, newData any) error {
	return vd.DB.Model(&visit).Where("id = ?", visit.ID).Updates(newData).Error
}

func (vd *VisitDB) UpdateStatus(ID int, status string) error {
	return vd.DB.Model(&model.Visit{}).Where("id = ?", ID).Update("status", status).Error
}

// GetVisitsByAdoperID get all visits by adopter id
// user is able to see all visits that they have
// return a list of visits that the adopter has made to the pets
func (vd *VisitDB) GetVisitsByAdoperID(adopterID uint) ([]model.Visit, error) {
	var visits []model.Visit
	err := vd.DB.Where("user_id = ?", adopterID).Preload("Pet", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Owner.Details")
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
	err := vd.DB.Where("owner_pet_id = ? and status = 'pending'", ownerID).Preload("User.Details").Preload("Pet").Find(&visits).Error

	if err != nil {
		return nil, err
	}

	return visits, nil
}
