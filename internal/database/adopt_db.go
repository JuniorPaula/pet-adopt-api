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
