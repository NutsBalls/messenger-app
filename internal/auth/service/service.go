package service

import (
	"github.com/google/uuid"
)

func createUser(name string, email string, password string) User {
	var user User

	user.ID = uuid.New()
	user.Name = name
	user.Email = email
	user.Password = password

	return user
}
