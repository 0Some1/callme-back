package models

import "gorm.io/gorm"

type Request struct {
	gorm.Model
	UserID   uint
	Follower User `gorm:"foreignKey:UserID"`
	Accepted bool `gorm:"type:bool;default:false"`
}
