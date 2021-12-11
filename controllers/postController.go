package controllers

import (
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

	if user.IsFollowing(otherUser.ID) {
		return c.JSON(otherUser.Posts)
	}

	otherUser.RemovePrivatePosts()

	return c.JSON(otherUser.Posts)

}

func GetExplorePosts(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	//pagination
	resultsPerPage, err := strconv.Atoi(c.Params("resultsPerPage"))
	page, err := strconv.Atoi(c.Params("page"))
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
	return c.JSON(posts)
}
