package model

import (
	"errors"
	"time"
)

type Adoption struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	PetID      uint      `json:"pet_id" gorm:"unique;not null"`
	OldOwnerID *uint     `json:"old_owner_id"`
	AdopterID  uint      `json:"adopter_id" gorm:"not null"`
	AdoptDate  time.Time `json:"adopt_date" gorm:"default:current_timestamp"`

	Pet      Pet   `json:"pet" gorm:"foreignKey:PetID"`
	OldOwner *User `json:"old_owner,omitempty" gorm:"foreignKey:OldOwnerID"`
	Adopter  User  `json:"adopter" gorm:"foreignKey:AdopterID"`
}

func NewAdoption(petID, oldOwnerID, adopterID uint) *Adoption {
	return &Adoption{
		PetID:      petID,
		OldOwnerID: &oldOwnerID,
		AdopterID:  adopterID,
		AdoptDate:  time.Now(),
	}
}

func (a *Adoption) ValidateFields() error {
	if a.PetID <= 0 {
		return errors.New("pet_id is required")
	}
	if a.OldOwnerID != nil {
		return errors.New("old_owner_id is required")
	}
	if a.AdopterID <= 0 {
		return errors.New("adopter_id is required")
	}
	return nil
}
