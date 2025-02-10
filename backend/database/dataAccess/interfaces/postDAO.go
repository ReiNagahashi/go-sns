package interfaces

import "go-sns/models"

type PostDAO interface{
	AddFavorite(userId, postId int) error
	Create(postData models.Post) error
	Delete(id int) error
	DeleteFavorite(userId, postId int) error
	GetAll() ([]models.Post, error)
	GetById(id int) (*models.Post, error)
	GetPosts(limit int)([]models.Post, error)
	ValidatePostField(postData models.Post) error
}