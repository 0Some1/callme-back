package DTO

import "callme/models"

type PostDTO struct {
	ID          uint
	UserID      uint
	UserName    string
	Avatar      string
	Bio         string
	Title       string
	Photos      []*models.Photo
	Description string
	Keywords    string
	Likes       int
	HasLiked    bool
	Comments    []*CommentDTO
}
