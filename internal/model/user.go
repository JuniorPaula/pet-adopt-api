package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"-"`
	IsAdmin   bool   `json:"is_admin"`
}

func NewUser(firstName, lastName, email, password string, isAdmin bool) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(hash),
		IsAdmin:   isAdmin,
	}, nil
}
