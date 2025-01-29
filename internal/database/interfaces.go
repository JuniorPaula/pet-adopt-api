package database

import "get_pet/internal/model"

type UserInterface interface {
	Create(user *model.User) error
}
