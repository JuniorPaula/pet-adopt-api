package database

import (
	"get_pet/internal/model"

	"gorm.io/gorm"
)

type UserDetailsDB struct {
	DB *gorm.DB
}

func NewUserDetails(db *gorm.DB) *UserDetailsDB {
	return &UserDetailsDB{
		DB: db,
	}
}

func (u *UserDetailsDB) Update(userDetails *model.UserDetails) error {
	return u.DB.Model(&userDetails).Where("user_id = ?", userDetails.UserID).Updates(userDetails).Error
}
