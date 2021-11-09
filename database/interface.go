package database

import "callme/models"

type databaseInterface interface {
	CreateUser(user *models.User) error
	SaveUser(user *models.User) error
	DeleteUser(user *models.User) error
	GetUserByID(userID string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	PreloadFollowers(user *models.User) error
	PreloadFollowings(user *models.User) error
	PreloadPosts(user *models.User) error
	CreatePost(post *models.Post) error
	CreatePhoto(photo *models.Photo) error
	GetPostByID(postID string) (*models.Post, error)
	GetPostByPhotoName(photoName string) (*models.Post, error)
}
