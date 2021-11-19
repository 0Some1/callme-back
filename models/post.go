package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID      uint
	Private     *bool     `json:"private,omitempty" gorm:"type:bool;default:false"`
	Title       string    `json:"title,omitempty" validate:"required"`
	Photos      []Photo   `json:"photos,omitempty" gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Description string    `json:"description,omitempty" validate:"required"`
	Keywords    string    `json:"keywords,omitempty" validate:"required"`
	Likes       []User    `json:"likes,omitempty" gorm:"many2many:user_like;"`
	Comments    []Comment `json:"comments,omitempty"`
}

func (p *Post) PreparePost(baseURL string) {
	for i := range p.Photos {
		p.Photos[i].AddBaseURL(baseURL)
	}
}
