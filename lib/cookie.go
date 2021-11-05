package lib

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

func SetCookie(id uint) (fiber.Cookie, error) {
	token, err := GenerateJwt(strconv.Itoa(int(id)))
	if err != nil {
		fmt.Print("lib.SetCookie - ")
		return fiber.Cookie{}, err
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: false,
		SameSite: "None",
		Secure:   false,
		Expires:  time.Now().Add(time.Hour * 24 * 9999),
	}

	return cookie, nil
}
