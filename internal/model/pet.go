package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Pet struct {
	gorm.Model
	ID        int      `json:"id"`
	UserID    uint     `json:"user_id"`
	Name      string   `json:"name" gorm:"type:varchar(50);not null"`
	Age       int      `json:"age" gorm:"not null"`
	Weight    float64  `json:"weight" gorm:"not null"`
	Size      string   `json:"size" gorm:"type:varchar(20)"`
	Color     string   `json:"color" gorm:"type:varchar(30)"`
	Images    []string `json:"images" gorm:"type:text[]"`
	Available bool     `json:"available" gorm:"default:true"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Owner User `json:"owner" gorm:"foreignKey:UserID"`
}

func NewPet(userID, age int, weight float64, name, size, color string, images []string) *Pet {
	return &Pet{
		UserID: uint(userID),
		Name:   name,
		Age:    age,
		Weight: weight,
		Size:   size,
		Color:  color,
		Images: images,
	}
}

func (p *Pet) ValidateFields() error {
	if p.Name == "" {
		return errors.New("pet name is required")
	}
	if p.Age <= 0 {
		return errors.New("pet age is required")
	}
	if p.Weight <= 0 {
		return errors.New("pet wight is required")
	}
	return nil
}
