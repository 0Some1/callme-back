package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID      uint
	Private     *bool     `json:"private,omitempty" gorm:"type:bool;default:false"`
	Title       string    `json:"title,omitempty"`
	Photos      []Photo   `json:"photos,omitempty" `
	Description string    `json:"description,omitempty"`
	Likes       []User    `gorm:"many2many:user_like;"`
	Comments    []Comment `json:"comments,omitempty"`
}
