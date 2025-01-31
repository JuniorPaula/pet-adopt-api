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
