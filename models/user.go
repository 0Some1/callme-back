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
	Post       []*Post    `json:"post,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Request    []*Request `json:"request,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Followers  []*User    `gorm:"many2many:user_follower;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"followers"`
	Followings []*User    `gorm:"many2many:user_following;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"followings"`
}
