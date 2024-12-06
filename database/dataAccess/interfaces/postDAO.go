package interfaces

import "go-sns/models"

type PostDAO interface{
	Create(postData models.Post) bool
	GetById(id int) models.Post
	Delete(id int) bool
	GetAll() []models.Post
}