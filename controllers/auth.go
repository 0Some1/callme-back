package controllers

import (
	"callme/database"
	"callme/models"
	"callme/utilities"
	"callme/utilities/validators"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)
	if err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Password:  password,
		Username:  data["username"],
		Email:     data["email"],
	}

	errors := validators.ValidateStruct(user)
	if errors != nil {
		return c.JSON(errors)
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not sign up this email",
		})
	}
	cookie := setCookie(user)
	c.Cookie(&cookie)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not set the cookie!",
		})
	}

	return c.JSON(fiber.Map{
		"username": user.Username,
		"email":    user.Email,
		"token":    cookie.Value,
	})

}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)
	if err != nil {
		return err
	}

	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "user could not be found!",
		})
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"]))
	if err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "incorrect password!",
		})
	}

	cookie := setCookie(user)
	c.Cookie(&cookie)

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not set the cookie!",
		})
	}

	return c.JSON(fiber.Map{
		"username": user.Username,
		"email":    user.Email,
		"token":    cookie.Value,
	})
}

func AuthenticatedUser(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	id, _ := utilities.ParseJwt(cookie)

	var user models.User

	database.DB.Where("id = ?", id).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: false,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func setCookie(user models.User) fiber.Cookie {
	token, _ := utilities.GenerateJwt(strconv.Itoa(int(user.Id)))

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: false,
		SameSite: "None",
		Secure:   false,
	}

	return cookie
}
