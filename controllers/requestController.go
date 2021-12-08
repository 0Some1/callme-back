package controllers

import (
	"callme/database"
	"callme/models"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func GetRequests(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	requests, err := database.DB.GetRequests(strconv.Itoa(int(user.ID)))
	if err != nil {
		fmt.Println("getRequestController - getReqsDB - ", err)
		return fiber.ErrInternalServerError
	}
	for i := 0; i < len(requests); i++ {
		requests[i].Follower.PrepareUser(c.BaseURL())
	}
	return c.JSON(requests)

}

func CreateRequest(c *fiber.Ctx) error {
	followingUserID := c.Params("id")
	user := c.Locals("user").(*models.User)
	if strconv.Itoa(int(user.ID)) == followingUserID {
		return fiber.NewError(fiber.StatusBadRequest, "You can't follow yourself")
	}
	request, err := database.DB.CreateRequest(user.ID, followingUserID)
	if err != nil {
		fmt.Println("createRequestController - createReqDB - ", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusBadRequest, "User not found")
		}
		if strings.Contains(err.Error(), "already exists") {
			return fiber.NewError(fiber.StatusBadRequest, "User already followed")
		}
		return fiber.ErrInternalServerError
	}
	return c.Status(201).JSON(request)
}

func DeleteRequest(c *fiber.Ctx) error {
	otherUserID := c.Params("userID")
	user := c.Locals("user").(*models.User)
	rowsAffected, err := database.DB.DeleteRequest(user.ID, otherUserID)
	if err != nil {
		fmt.Println("deleteRequestController - deleteReqDB - ", err)
		if rowsAffected == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "Request not found")
		}
		return fiber.ErrInternalServerError
	}
	return c.Status(204).JSON(nil)
}

func AcceptRequest(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	requestID := c.Params("id")
	err := database.DB.AcceptRequest(requestID, user)
	if err != nil {
		fmt.Println("acceptRequestController - acceptReqDB - ", err)
		return fiber.ErrBadRequest
	}
	return c.Status(204).JSON(nil)
}

func DeclineRequest(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	requestID := c.Params("id")
	err := database.DB.DeclineRequest(requestID, user)
	if err != nil {
		fmt.Println("acceptRequestController - acceptReqDB - ", err)
		return fiber.ErrBadRequest
	}
	return c.Status(204).JSON(nil)
}
