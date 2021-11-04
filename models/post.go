package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID      uint
	Title       string `json:"title"`
	Photos      string `json:"photos"`
	Description string `json:"description"`
	Likes       []User `gorm:"many2many:user_like;"`
	Comments    []Comment
}
