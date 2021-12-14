package DTO

import (
	"callme/models"
)

type CommentDTO struct {
	ID         uint
	UserID     uint
	PostID     uint
	Text       string
	OwnComment bool //if the requesting user has posted this comment
}

func PrepareCommentDTO(userID uint, comments []*models.Comment) []*CommentDTO {
	commentDTOs := make([]*CommentDTO, 0)
	for i := range comments {
		commentDTOs = append(commentDTOs, &CommentDTO{
			ID:         comments[i].ID,
			UserID:     comments[i].UserID,
			PostID:     comments[i].PostID,
			Text:       comments[i].Text,
			OwnComment: comments[i].UserID == userID,
		})
	}
	return commentDTOs
}
