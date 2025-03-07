package interfaces

import "go-sns/models"

type PostDAO interface{
	Create(postData models.Post) error
	Delete(id int) error
	GetAll(...int) ([]models.Post, error)
	GetById(id int) (*models.Post, error)
	ValidatePostField(postData models.Post) error
}