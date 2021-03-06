package database

import (
	"callme/models"
	"errors"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type postgresDB struct {
	DatabaseInterface
	db *gorm.DB
}

func (p *postgresDB) CreateUser(user *models.User) error {
	return p.db.Create(user).Error
}

func (p postgresDB) SaveUser(user *models.User) error {
	return p.db.Save(user).Error
}
func (p postgresDB) DeleteUser(user *models.User) error {
	return p.db.Delete(user).Error
}

func (p *postgresDB) GetUserByID(userID string) (*models.User, error) {
	user := new(models.User)
	err := p.db.Where("id = ?", userID).First(&user).Error
	return user, err
}

func (p *postgresDB) SearchUsers(search string) ([]*models.User, error) {
	users := make([]*models.User, 0)
	search = strings.ToLower(search)
	err := p.db.Where("lower(name) LIKE ?", "%"+search+"%").Or("lower(username) LIKE ?", "%"+search+"%").Find(&users).Error
	return users, err
}

func (p *postgresDB) GetUserByEmail(email string) (*models.User, error) {
	user := new(models.User)
	err := p.db.Where("email = ?", email).First(&user).Error
	return user, err
}
func (p *postgresDB) PreloadFollowers(user *models.User) error {
	err := p.db.Preload("Followers").Where("id = ?", user.ID).Find(&user).Error
	return err
}

func (p *postgresDB) PreloadFollowings(user *models.User) error {
	err := p.db.Preload("Followings").Where("id = ?", user.ID).Find(&user).Error
	return err
}

func (p *postgresDB) PreloadPosts(user *models.User) error {
	err := p.db.Preload("Posts.Photos").Where("id = ?", user.ID).Find(&user).Error
	return err
}

func (p *postgresDB) LoadExplorePosts(user *models.User, resultsPerPage int, page int) ([]*models.Post, error) {
	posts := make([]*models.Post, 0)
	err := p.db.Preload("Photos").
		Where("user_id != ?", user.ID).
		Where("private = false OR user_id IN (SELECT following_id AS user_id FROM user_following WHERE user_id = ?)", user.ID).
		Order("created_at DESC").
		Offset(resultsPerPage * (page - 1)).
		Limit(resultsPerPage).
		Find(&posts).Error
	return posts, err
}

func (p *postgresDB) AddCommentToPost(comment *models.Comment) error {
	err := p.db.Model(&models.Comment{}).Create(comment).Error
	return err
}

func (p *postgresDB) DeleteComment(commentID string, userID string) error {
	err := p.db.Where("id = ? AND user_id = ?", commentID, userID).Delete(&models.Comment{}).Error
	return err
}

func (p *postgresDB) PreloadRequests(user *models.User) error {
	err := p.db.Preload("Requests.Follower").Where("id = ?", user.ID).Find(&user).Error
	return err
}

func (p *postgresDB) CreatePost(post *models.Post) error {
	return p.db.Create(post).Error
}

func (p *postgresDB) EditPost(postID uint, post *models.Post) error {
	return p.db.Model(&models.Post{}).Where("id = ?", postID).Save(&post).Error
}

func (p *postgresDB) DeletePost(post *models.Post) (int64, error) {
	err := p.db.Unscoped().Select(clause.Associations).Delete(&post)
	return err.RowsAffected, err.Error
}

func (p *postgresDB) CreatePhoto(photo *models.Photo) error {
	return p.db.Create(photo).Error
}

func (p *postgresDB) GetPostByID(postID string) (*models.Post, error) {
	post := new(models.Post)
	err := p.db.Where("id = ?", postID).First(&post).Error
	return post, err
}

func (p *postgresDB) PreloadPostByID(postID string) (*models.Post, error) {
	post := new(models.Post)
	err := p.db.Preload("Photos").Preload("Comments").Preload("Likes").Where("id = ?", postID).
		First(&post).Error
	return post, err
}

func (p *postgresDB) GetPostByPhotoName(photoName string) (*models.Post, error) {
	//I didn't read the whole doc of gorm, but it must be a better way to do this
	photo := new(models.Photo)
	err := p.db.Model(&models.Photo{}).Where("name = ? ", photoName).Find(&photo).Error
	post, err := p.GetPostByID(strconv.Itoa(int(photo.PostID)))
	return post, err
}

func (p *postgresDB) LikePost(postID string, userID string) error {
	err := p.db.Table("user_like").Create([]map[string]interface{}{
		{"post_id": postID, "user_id": userID},
	}).Error
	return err
}

func (p *postgresDB) UnlikePost(postID string, userID string) error {
	err := p.db.Table("user_like").Where("user_id = ? AND post_id = ?", userID, postID).Unscoped().Delete([]map[string]interface{}{
		{"post_id": postID, "user_id": userID},
	})

	if err.RowsAffected == 0 {
		return errors.New("This post was not liked by this user")
	}

	return err.Error
}

func (p *postgresDB) GetRequests(id string) ([]*models.Request, error) {
	requests := make([]*models.Request, 0)
	err := p.db.Where("user_id = ?", id).Preload("Follower").Find(&requests).Error
	return requests, err
}

// GetRequestByID returns the request that created before from userID to requestUserID
func (p *postgresDB) GetRequestByID(userID uint, requestUserID string) (*models.Request, int) {
	request := new(models.Request)
	err := p.db.Where("user_id = ? AND follower_id = ?", requestUserID, userID).First(&request).RowsAffected
	return request, int(err)
}
func (p *postgresDB) CreateRequest(userID uint, requestUserID string) (*models.Request, error) {
	request := new(models.Request)
	//first must check if the request already exists
	request, rowsAffected := p.GetRequestByID(userID, requestUserID)
	if rowsAffected != 0 {
		return nil, errors.New("request already exists")
	}
	//and if the user the requested user exists too
	requestedUser, err := p.GetUserByID(requestUserID)
	if err != nil {
		return nil, err
	}
	//then create the request
	request.FollowerID = userID
	requestedUser.Requests = append(requestedUser.Requests, request)
	err = p.db.Save(requestedUser).Error
	return request, err
}

func (p *postgresDB) DeleteRequest(userID uint, otherUserID string) (int64, error) {
	request, rows := p.GetRequestByID(userID, otherUserID)
	if rows == 0 {
		return 0, errors.New("request not found")
	}
	err := p.db.Unscoped().Delete(request)
	return err.RowsAffected, err.Error
}

func (p *postgresDB) PreLoadPhotos(post *models.Post) error {
	err := p.db.Preload("Photos").Where("id = ?", post.ID).Find(&post).Error
	return err
}

func (p *postgresDB) IsFollowing(ownUserID uint, otherUserID uint) (error, bool) {
	isFollowing := int64(0)
	err := p.db.Table("user_following").Where("user_id = ? AND following_id = ?", ownUserID, otherUserID).Count(&isFollowing).Error
	return err, isFollowing != 0
}

func (p *postgresDB) AcceptRequest(requestID string, user *models.User) error {
	request := new(models.Request)
	err := p.db.Where("id = ?", requestID).First(&request).Error
	if err != nil {
		return err
	}
	if user.ID != request.UserID {
		return errors.New("you are not the owner of this request")
	}
	err = p.FollowByID(request.FollowerID, user.ID)
	if err != nil {
		return err
	}
	err = p.db.Unscoped().Delete(&request).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *postgresDB) DeclineRequest(requestID string, user *models.User) error {
	request := new(models.Request)
	err := p.db.Where("id = ?", requestID).First(&request).Error
	if err != nil {
		return err
	}
	if user.ID != request.UserID {
		return errors.New("you are not the owner of this request")
	}
	err = p.db.Unscoped().Delete(&request).Error
	if err != nil {
		return err
	}
	return nil
}

// FollowByID userID is the user who is following otherUserID
func (p *postgresDB) FollowByID(userID uint, otherUserID uint) error {
	user, err := p.GetUserByID(strconv.Itoa(int(userID)))
	if err != nil {
		return err
	}
	otherUser, err := p.GetUserByID(strconv.Itoa(int(otherUserID)))
	if err != nil {
		return err
	}
	p.db.Table("user_following").Create([]map[string]interface{}{
		{"user_id": user.ID, "following_id": otherUser.ID},
	})
	p.db.Table("user_follower").Create([]map[string]interface{}{
		{"user_id": otherUser.ID, "follower_id": user.ID},
	})

	return nil
}
func (p *postgresDB) Unfollow(user *models.User, otherUser *models.User) error {

	err := p.db.Table("user_following").Where("user_id = ? AND following_id = ?", user.ID, otherUser.ID).Unscoped().Delete([]map[string]interface{}{
		{"user_id": user.ID, "following_id": otherUser.ID},
	})
	err2 := p.db.Table("user_follower").Where("user_id = ? AND follower_id = ?", otherUser.ID, user.ID).Unscoped().Delete([]map[string]interface{}{
		{"user_id": otherUser.ID, "follower_id": user.ID},
	})
	if err.Error != nil || err2.Error != nil || err.RowsAffected == 0 || err2.RowsAffected == 0 {
		return errors.New("error unfollowing! maybe you are not following this user")
	}
	return nil
}
