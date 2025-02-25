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

func (vd *VisitDB) GetByPetID(petID int) (*model.Visit, error) {
	var v *model.Visit
	err := vd.DB.Where("pet_id = ?", petID).First(&v).Error
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

func (vd *VisitDB) GetVisitsByUserID(userID uint) ([]model.Visit, error) {
	var visits []model.Visit
	err := vd.DB.Where("user_id = ?", userID).Preload("Pet").Find(&visits).Error
	if err != nil {
		return nil, err
	}

	return visits, nil
}
