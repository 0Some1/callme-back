package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username   string     `json:"username,omitempty" gorm:"unique" validate:"required,min=3"`
	Name       string     `json:"name,omitempty"`
	Email      string     `json:"email,omitempty" gorm:"unique" validate:"required,email"`
	Password   string     `json:"password,omitempty" validate:"required,min=8"`
	Country    string     `json:"country"  validate:"required"`
	City       string     `json:"city" validate:"required"`
	Born       *time.Time `json:"born" validate:"required"`
	Avatar     string     `json:"avatar" validate:"required"`
	Bio        string     `json:"bio" validate:"required"`
	Posts      []*Post    `json:"post,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Requests   []*Request `json:"request,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Followers  []*User    `gorm:"many2many:user_follower;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"followers"`
	Followings []*User    `gorm:"many2many:user_following;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"followings"`
}

func (u User) IsFollowing(id uint) bool {
	for _, user := range u.Followings {
		if user.ID == id {
			return true
		}
	}
	return false
}

func (u *User) PrepareUser(baseurl string) {
	u.Password = ""
	u.Followers = nil
	u.Followings = nil
	if u.Avatar != "" {
		u.Avatar = baseurl + u.Avatar
	}
}

func (u *User) RemovePrivatePosts() {
	for i, post := range u.Posts {
		if *post.Private {
			u.Posts = append(u.Posts[:i], u.Posts[i+1:]...)
		}
	}
}
