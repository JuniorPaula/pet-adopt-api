package database

import (
	"get_pet/internal/model"
	"time"

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
	result := u.DB.Create(user)

	if err := u.DB.Save(&model.UserDetails{
		UserID:    user.ID,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05")},
	).Error; err != nil {
		return err
	}

	return result.Error
}

func (u *UserDB) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := u.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (u *UserDB) GetByID(id int) (*model.User, error) {
	var user model.User
	err := u.DB.Preload("Details").Where("id = ?", id).First(&user).Error
	return &user, err
}

func (u *UserDB) Update(user *model.User, newUser any) error {
	return u.DB.Model(&user).Where("id = ?", user.ID).Updates(newUser).Error
}

func (u *UserDB) SoftRemove(id int) error {
	return u.DB.Model(&model.User{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}
