package database

import "callme/models"

type DatabaseInterface interface {
	CreateUser(user *models.User) error
	SaveUser(user *models.User) error
	DeleteUser(user *models.User) error
	GetUserByID(userID string) (*models.User, error)
	SearchUsers(search string) ([]*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	PreloadFollowers(user *models.User) error
	PreloadFollowings(user *models.User) error
	PreloadPosts(user *models.User) error
	PreloadRequests(user *models.User) error
	CreatePost(post *models.Post) error
	CreatePhoto(photo *models.Photo) error
	GetPostByID(postID string) (*models.Post, error)
	GetPostByPhotoName(photoName string) (*models.Post, error)
	GetRequests(id string) ([]*models.Request, error)
	GetRequestByID(userID uint, requestUserID string) (*models.Request, int)
	CreateRequest(userID uint, requestUserID string) (*models.Request, error)
	DeleteRequest(requestID string) error
	AcceptRequest(requestID string, user *models.User) error
	FollowByID(userID uint, otherUserID uint) error
}
