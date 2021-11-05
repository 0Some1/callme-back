package controllers

import (
	"callme/database"
	"callme/lib"
	"callme/models"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	var user models.User
	err := c.BodyParser(&user)
	if err != nil {
		fmt.Println("RegisterController - BodyParser - ", err)
		c.Status(fiber.ErrNotAcceptable.Code)
		return c.JSON(lib.CustomError(fiber.ErrNotAcceptable, "can't read body as JSON!"))
	}

	validationErrors := lib.ValidateRegister(user)
	if validationErrors != nil {
		return c.Status(fiber.ErrForbidden.Code).JSON(validationErrors)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		fmt.Println("RegisterController - GenerateFromPassword - ", err)
		c.Status(fiber.ErrInternalServerError.Code)
		return c.JSON(lib.CustomError(fiber.ErrInternalServerError, ""))
	}
	user.Password = string(password)

	err = database.DB.CreateUser(&user)
	if err != nil {
		fmt.Println("RegisterController - createUser - ", err)
		c.Status(fiber.ErrNotAcceptable.Code)
		return c.JSON(lib.CustomError(fiber.ErrNotAcceptable, "the email must be unique!"))
	}

	cookie, err := lib.SetCookie(user.ID)
	if err != nil {
		fmt.Println("RegisterController - setCookie - ", err)
		c.Status(fiber.ErrInternalServerError.Code)
		return c.JSON(lib.CustomError(fiber.ErrInternalServerError, ""))
	}
	c.Cookie(&cookie)

	return c.Status(201).JSON(fiber.Map{
		"token": cookie.Value,
	})

}

func Login(c *fiber.Ctx) error {
	var userIn models.User
	err := c.BodyParser(&userIn)
	if err != nil {
		fmt.Println("LoginController - BodyParser - ", err)
		c.Status(fiber.ErrNotAcceptable.Code)
		return c.JSON(lib.CustomError(fiber.ErrNotAcceptable, "can't read body as JSON!"))
	}

	user, err := database.DB.GetUserByEmail(userIn.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("LoginController - First - ", err)
			c.Status(fiber.ErrNotFound.Code)
			return c.JSON(lib.CustomError(fiber.ErrNotFound, "user not found!"))
		}
		fmt.Println("LoginController - First - ", err)
		c.Status(fiber.ErrInternalServerError.Code)
		return c.JSON(lib.CustomError(fiber.ErrInternalServerError, ""))
	}

	if user.ID == 0 {
		c.Status(fiber.ErrNotFound.Code)
		return c.JSON(lib.CustomError(fiber.ErrNotFound, "User not found!"))
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userIn.Password))
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(lib.CustomError(fiber.ErrUnauthorized, "Password is incorrect!"))
	}

	cookie, err := lib.SetCookie(user.ID)
	if err != nil {
		fmt.Println("RegisterController - setCookie - ", err)
		c.Status(fiber.ErrInternalServerError.Code)
		return c.JSON(lib.CustomError(fiber.ErrInternalServerError, ""))
	}
	c.Cookie(&cookie)

	return c.Status(200).JSON(fiber.Map{
		"token": cookie.Value,
	})
}
