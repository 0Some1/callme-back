package database

import (
	"callme/models"
	"errors"
	"gorm.io/gorm"
	"strconv"
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

func (p *postgresDB) CreatePost(post *models.Post) error {
	return p.db.Create(post).Error
}

func (p *postgresDB) CreatePhoto(photo *models.Photo) error {
	return p.db.Create(photo).Error
}

func (p *postgresDB) GetPostByID(postID string) (*models.Post, error) {
	post := new(models.Post)
	err := p.db.Where("id = ?", postID).First(&post).Error
	return post, err
}

func (p *postgresDB) GetPostByPhotoName(photoName string) (*models.Post, error) {
	//I didn't read the whole doc of gorm, but it must be a better way to do this
	photo := new(models.Photo)
	err := p.db.Model(&models.Photo{}).Where("name = ? ", photoName).Find(&photo).Error
	post, err := p.GetPostByID(strconv.Itoa(int(photo.PostID)))
	return post, err
}
func (p *postgresDB) GetRequests(id string) ([]*models.Request, error) {
	requests := make([]*models.Request, 0)
	err := p.db.Where("user_id = ?", id).Preload("Follower").Find(&requests).Error
	return requests, err
}
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

func (p *postgresDB) DeleteRequest(reqID string) error {
	err := p.db.Unscoped().Delete(&models.Request{}, reqID).Error
	return err
}

func (p *postgresDB) PreLoadPhotos(post *models.Post) error {
	err := p.db.Preload("Photos").Where("id = ?", post.ID).Find(&post).Error
	return err
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
	err = p.FollowByID(request.Follower.ID, user.ID)
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
	user.Followings = append(user.Followings, otherUser)
	otherUser.Followers = append(otherUser.Followers, user)
	err = p.db.Save(user).Error
	if err != nil {
		return err
	}
	err = p.db.Save(otherUser).Error
	return err
}
