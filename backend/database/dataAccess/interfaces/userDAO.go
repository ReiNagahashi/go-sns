package interfaces

import "go-sns/models"

type UserDAO interface{
	Create(userData *models.User, password string) error
	GetAll(...int)([]models.User, error)
	GetById(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetHashedPasswordById(id int) string
	ValidateUserField(string, string, bool) error
}