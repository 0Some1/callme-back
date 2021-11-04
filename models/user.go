package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string `validate:"required,min=3,max=32" json:"username,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `gorm:"unique" validate:"required,email,max=32" json:"email"`
	Password    string `validate:"required" json:"password"`
	Post        []Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	followerID  *uint
	followingID *uint
	followers   []User `gorm:"foreignkey:followerID"`
	followings  []User `gorm:"foreignkey:followingID"`
}
