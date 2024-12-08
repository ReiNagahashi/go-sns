package interfaces

import "go-sns/models"

type UserDAO interface{
	Create(userData models.User) bool
	GetById(id int) []models.User
	GetByEmail(email string) models.User
	GetHashedPasswordById(id int) string
}