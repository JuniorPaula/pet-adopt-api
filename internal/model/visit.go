package model

import (
	"errors"
	"time"
)

type Visit struct {
	ID     int       `json:"id" gorm:"primaryKey"`
	UserID uint      `json:"user_id" gorm:"not null"`
	PetID  uint      `json:"pet_id" gorm:"not null"`
	Date   time.Time `json:"date" gorm:"not null"`
	Status string    `json:"status" gorm:"default:'pending'"`

	User User `json:"user" gorm:"foreignKey:UserID"`
	Pet  Pet  `json:"pet" gorm:"foreignKey:PetID"`
}

func NewVisit(userId, petID int, status string) *Visit {
	return &Visit{
		UserID: uint(userId),
		PetID:  uint(petID),
		Status: status,
		Date:   time.Now(),
	}
}

func (v *Visit) ValidateFields() error {
	if v.UserID <= 0 {
		return errors.New("user_id is required")
	}
	if v.PetID <= 0 {
		return errors.New("pet_id is required")
	}
	return nil
}
