package controllers

import (
	"callme/DTO"
	"callme/database"
	"callme/lib"
	"callme/models"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreatePost(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	post := new(models.Post)
	err := c.BodyParser(post)
	if err != nil {
		fmt.Println("CreatePost - Bodyparser - ", err)
		return fiber.ErrBadRequest
	}
	//post validation must be done here!!!
	err = validator.New().StructPartial(post, "Title", "Description")
	if err != nil {
		fmt.Println("CreatePost - Validation - ", err)
		return fiber.ErrBadRequest
	}

	post.UserID = user.ID

	var photos []*models.Photo

	if form, err := c.MultipartForm(); err == nil {
		files := form.File["photos"]
		if files == nil {
			return fiber.NewError(fiber.StatusBadRequest, "No photos were uploaded")
		}
		for _, file := range files {
			err = lib.ImageValidation(file)
			if err != nil {
				fmt.Println("CreatePost - imageValidation -", err)
				return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
			}
			file.Filename = lib.GenFileName(file.Filename)
			err = c.SaveFile(file, fmt.Sprintf("./uploads/post/%s", file.Filename))
			if err != nil {
				fmt.Println("CreatePost - saveFile ", err)
				return fiber.ErrInternalServerError
			}
			photoTemp := models.Photo{
				Name: file.Filename,
				Path: "/uploads/post/" + file.Filename,
			}
			photos = append(photos, &photoTemp)
			err := database.DB.CreatePhoto(&photoTemp)
			if err != nil {
				fmt.Println("CreatePost - createPhoto -", err)
				return fiber.ErrInternalServerError
			}

		}

	}

	for i := 0; i < len(photos); i++ {
		photos[i].Path = c.BaseURL() + photos[i].Path
	}
	post.Photos = photos
	err = database.DB.CreatePost(post)
	if err != nil {
		fmt.Println("CreatePost - SavePost -", err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(post)
}

func EditPost(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	postID := c.Params("postID")
	newPost := new(models.Post)
	err := c.BodyParser(newPost)

	if err != nil {
		fmt.Println("EditPost - Bodyparser - ", err)
		return fiber.ErrBadRequest
	}
	//post validation must be done here!!!
	// err = validator.New().StructPartial(newPost, "Title", "Description")
	// if err != nil {
	// 	fmt.Println("EditPost - Validation - ", err)
	// 	return fiber.ErrBadRequest
	// }

	oldPost, err := database.DB.GetPostByID(postID)
	if err != nil {
		fmt.Println("EditPost - Get old post - ", err)
		return fiber.ErrInternalServerError
	}

	//check if user posted this post
	if oldPost.UserID != user.ID {
		return fiber.NewError(fiber.StatusForbidden, "This user is not the owner of the post.")
	}

	if newPost.Title != "" {
		oldPost.Title = newPost.Title
	}
	if newPost.Description != "" {
		oldPost.Description = newPost.Description
	}
	if newPost.Keywords != "" {
		oldPost.Keywords = newPost.Keywords
	}
	if newPost.Private != nil {
		println(newPost.Private)
		oldPost.Private = newPost.Private
	}
	err = database.DB.EditPost(oldPost.ID, oldPost)
	if err != nil {
		fmt.Println("EditPost - EditPost -", err)
		return fiber.ErrInternalServerError
	}

	return c.Status(204).JSON(nil)

}

func DeletePost(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	postID := c.Params("postID")

	//get post
	post, err := database.DB.GetPostByID(postID)
	//check if post exists
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Post not found")
	}
	//check if user posted this post
	if post.UserID != user.ID {
		return fiber.NewError(fiber.StatusForbidden, "Can not delete this post")
	}

	//delete post from database
	rowsAffected, err := database.DB.DeletePost(post)
	if err != nil {
		fmt.Println("deletePostController - deletePostDB - ", err)
		if rowsAffected == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "Post not found")
		}
		return fiber.ErrInternalServerError
	}
	return c.Status(204).JSON(nil)
}

func GetPosts(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	err := database.DB.PreloadPosts(user)
	if err != nil {
		fmt.Println("GetPosts - ", err)
		return fiber.ErrInternalServerError
	}
	for _, post := range user.Posts {
		post.PreparePost(c.BaseURL())
	}
	return c.JSON(user.Posts)
}

func GetPostsByUserID(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	userID := c.Params("userID")
	otherUser, err := database.DB.GetUserByID(userID)
	if err != nil {
		fmt.Println("GetPostsByUserID - GetUserByID  ", err)
		return fiber.ErrNotFound
	}
	err = database.DB.PreloadFollowings(user)
	if err != nil {
		fmt.Println("GetPostsByUserID - PreloadFollowings  ", err)
		return fiber.ErrInternalServerError
	}
	err = database.DB.PreloadPosts(otherUser)
	if err != nil {
		fmt.Println("GetPostsByUserID - PreloadPosts  ", err)
		return fiber.ErrInternalServerError
	}

	for _, post := range otherUser.Posts {
		post.PreparePost(c.BaseURL())
	}

	_, isFollowing := database.DB.IsFollowing(user.ID, otherUser.ID)
	//check if user has access to the post
	if isFollowing {
		return c.JSON(otherUser.Posts)
	}

	otherUser.RemovePrivatePosts()

	return c.JSON(otherUser.Posts)

}

func GetPostDetails(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	postID := c.Params("postID")
	//get the post
	post, err := database.DB.PreloadPostByID(postID)
	if err != nil {
		fmt.Println("GetPostDetails - GetPostByID  ", err)
		return fiber.ErrNotFound
	}
	_, isFollowing := database.DB.IsFollowing(user.ID, post.UserID)
	//check if user has access to the post
	if !isFollowing && *post.Private && post.UserID != user.ID {
		return fiber.NewError(fiber.StatusForbidden, "Can not get this post")
	}

	post.PreparePost(c.BaseURL())

	comments := DTO.PrepareCommentDTOs(user.ID, post.Comments)
	hasLiked := user.HasLikedPost(post.Likes)
	postOwner, err := database.DB.GetUserByID(strconv.FormatUint(uint64(post.UserID), 10))
	if err != nil {
		fmt.Println("GetPostDetails - GetUserByID  ", err)
		return fiber.ErrInternalServerError
	}

	postDTO := DTO.PostDTO{
		ID:          post.ID,
		UserID:      post.UserID,
		UserName:    postOwner.Username,
		Avatar:      postOwner.Avatar,
		Bio:         postOwner.Bio,
		Title:       post.Title,
		Photos:      post.Photos,
		Description: post.Description,
		Keywords:    post.Keywords,
		Likes:       len(post.Likes),
		HasLiked:    hasLiked,
		Comments:    comments,
	}

	return c.JSON(postDTO)
}

func GetExplorePosts(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	//pagination
	resultsPerPage, _ := strconv.Atoi(c.Query("resultsPerPage"))
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		return fiber.ErrBadRequest
	}
	if page <= 0 || resultsPerPage <= 0 {
		return fiber.ErrBadRequest
	}

	posts, err := database.DB.LoadExplorePosts(user, resultsPerPage, page)
	if err != nil {
		fmt.Println("GetExplorePosts - PreloadPosts  ", err)
		return fiber.ErrInternalServerError
	}

	for _, post := range posts {
		post.PreparePost(c.BaseURL())
	}

	return c.JSON(posts)
}

//set a comment on a post
func SetComment(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	postID := c.Params("postID")

	//get the post
	post, err := database.DB.GetPostByID(postID)
	if err != nil {
		fmt.Println("SetComment - GetPost  ", err)
		return fiber.ErrNotFound
	}

	//parse the comment
	comment := new(models.Comment)
	err = c.BodyParser(comment)
	comment.UserID = user.ID
	comment.PostID = post.ID
	if err != nil {
		fmt.Println("SetComment - Bodyparser - ", err)
		return fiber.ErrBadRequest
	}

	_, isFollowing := database.DB.IsFollowing(user.ID, post.UserID)
	//check if user has access to the post
	if !isFollowing && *post.Private && post.UserID != user.ID {
		return fiber.NewError(fiber.StatusForbidden, "Can not get this post")
	}

	err = database.DB.AddCommentToPost(comment)
	if err != nil {
		fmt.Println("SetComment - AddCommentToPost - ", err)
		return fiber.ErrInternalServerError
	}

	commentDTO := DTO.PrepareCommentDTO(user.ID, comment)
	return c.JSON(commentDTO)
}

//delete a comment from a post
func DeleteComment(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	commentID := c.Params("commentID")

	//delete the post
	err := database.DB.DeleteComment(commentID, strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		fmt.Println("DeleteComment - DeleteComment  ", err)
		return fiber.ErrInternalServerError
	}

	return c.Status(204).JSON(nil)
}

func LikePost(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	postID := c.Params("postID")
	//like or unlike
	like, err := strconv.Atoi(c.Query("like"))
	if err != nil {
		fmt.Println("LikePost - LikeQuery  ", err)
		return fiber.ErrInternalServerError
	}

	if like > 0 {
		err = database.DB.LikePost(postID, strconv.FormatUint(uint64(user.ID), 10))
	} else {
		err = database.DB.UnlikePost(postID, strconv.FormatUint(uint64(user.ID), 10))
	}

	if err != nil {
		fmt.Println("LikePost - LikePost  ", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(204).JSON(nil)
}
