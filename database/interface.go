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
	DeletePost(post *models.Post) (int64, error)
	CreatePhoto(photo *models.Photo) error
	GetPostByID(postID string) (*models.Post, error)
	PreloadPostByID(postID string) (*models.Post, error)
	GetPostByPhotoName(photoName string) (*models.Post, error)
	LoadExplorePosts(user *models.User, resultsPerPage int, page int) ([]*models.Post, error)
	AddCommentToPost(comment *models.Comment) error
	DeleteComment(commentID string) error
	LikePost(postID string, userID string) error
	GetRequests(id string) ([]*models.Request, error)
	GetRequestByID(userID uint, requestUserID string) (*models.Request, int)
	CreateRequest(userID uint, requestUserID string) (*models.Request, error)
	DeleteRequest(userID uint, otherUserID string) (int64, error)
	AcceptRequest(requestID string, user *models.User) error
	DeclineRequest(requestID string, user *models.User) error
	FollowByID(userID uint, otherUserID uint) error
	Unfollow(user *models.User, otherUser *models.User) error
}
