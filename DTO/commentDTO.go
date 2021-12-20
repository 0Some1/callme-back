package DTO

import (
	"callme/database"
	"callme/models"
	"strconv"
)

type CommentDTO struct {
	ID         uint
	UserID     uint
	UserName   string
	Avatar     string
	Bio        string
	PostID     uint
	Text       string
	OwnComment bool //if the requesting user has posted this comment
}

func PrepareCommentDTOs(userID uint, comments []*models.Comment) []*CommentDTO {
	commentDTOs := make([]*CommentDTO, 0)
	for i := range comments {
		userIDString := strconv.FormatUint(uint64(comments[i].UserID), 10)
		user, err := database.DB.GetUserByID(userIDString)
		if err != nil {
			continue
		}
		commentDTOs = append(commentDTOs, &CommentDTO{
			ID:         comments[i].ID,
			UserID:     comments[i].UserID,
			PostID:     comments[i].PostID,
			UserName:   user.Username,
			Avatar:     user.Avatar,
			Bio:        user.Bio,
			Text:       comments[i].Text,
			OwnComment: comments[i].UserID == userID,
		})
	}
	return commentDTOs
}

func PrepareCommentDTO(userID uint, comment *models.Comment) CommentDTO {
	userIDString := strconv.FormatUint(uint64(comment.UserID), 10)
	user, err := database.DB.GetUserByID(userIDString)
	if err != nil {
		return CommentDTO{}
	}
	commentDTOs := CommentDTO{
		ID:         comment.ID,
		UserID:     comment.UserID,
		PostID:     comment.PostID,
		UserName:   user.Username,
		Avatar:     user.Avatar,
		Bio:        user.Bio,
		Text:       comment.Text,
		OwnComment: comment.UserID == userID,
	}
	return commentDTOs
}
