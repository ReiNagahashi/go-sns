package interfaces

import "go-sns/models"

type PostDAO interface{
	Create(postData models.Post) bool
	Delete(id int) bool
	GetAll(...int) []models.Post
	GetById(id int) []models.Post
	ValidatePostField(postData models.Post) error
}