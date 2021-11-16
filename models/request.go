package models

import "gorm.io/gorm"

type Request struct {
	gorm.Model
	UserID     uint `json:"-"`
	FollowerID uint `json:"-"`
	Follower   User ` json:"follower" gorm:"foreignKey:FollowerID"`
}
