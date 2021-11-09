package database

import (
	"callme/models"
	"gorm.io/gorm"
	"strconv"
)

type postgresDB struct {
	databaseInterface
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
	err := p.db.Preload("Posts").Where("id = ?", user.ID).Find(&user).Error
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
