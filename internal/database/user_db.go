package database

import (
	"get_pet/internal/model"

	"gorm.io/gorm"
)

type UserDB struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *UserDB {
	return &UserDB{
		DB: db,
	}
}

func (u *UserDB) Create(user *model.User) error {
	return u.DB.Create(user).Error
}

func (u *UserDB) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := u.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}
