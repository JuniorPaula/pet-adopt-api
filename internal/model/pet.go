package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type StringArray []string

type Pet struct {
	ID          int         `json:"id"`
	UserID      uint        `json:"user_id"`
	Name        string      `json:"name" gorm:"type:varchar(50);not null"`
	Age         string      `json:"age" gorm:"type:varchar(30); not null"`
	Weight      string      `json:"weight" gorm:"type:varchar(30); not null"`
	Size        string      `json:"size" gorm:"type:varchar(20)"`
	Color       string      `json:"color" gorm:"type:varchar(30)"`
	Images      StringArray `json:"images" gorm:"type:jsonb"`
	Available   bool        `json:"available" gorm:"default:true"`
	Description string      `json:"description" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Owner User `json:"owner" gorm:"foreignKey:UserID"`
}

func NewPet(userID int, age, weight, name, size, color string, images []string) *Pet {
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
	if p.Age == "" {
		return errors.New("pet age is required")
	}
	if p.Weight == "" {
		return errors.New("pet wight is required")
	}
	if p.Size == "" {
		return errors.New("pet size is required")
	}
	if p.Color == "" {
		return errors.New("pet color is required")
	}
	return nil
}

func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringArray")
	}

	return json.Unmarshal(bytes, s)
}

const defaultLimitDescription = 1000

func (p *Pet) GetlimitDescription() int {
	return defaultLimitDescription
}
