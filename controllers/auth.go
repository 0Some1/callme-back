package controllers

import (
	"callme/database"
	"callme/models"
	"callme/utilities/validators"
	"github.com/dgrijalva/jwt-go"
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

	database.DB.Create(&user)
	return c.JSON(user)
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
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.StandardClaims{
		Issuer: strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // it will be last for 1 day
	})

	token , err := claims.SignedString([]byte("naughty"))
	if err!=nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	cookie := fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(user)
}

type Claims struct {
	jwt.StandardClaims
}
func AuthenticatedUser(c *fiber.Ctx) error  {
	cookie:=c.Cookies("jwt")
	token,err := jwt.ParseWithClaims(cookie,&Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("naughty"),nil
	})

	if err!=nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message":"unauthenticated!",
		})
	}
	claims := token.Claims.(*Claims)

	var user models.User
	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name: "jwt",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}