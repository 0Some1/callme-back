package lib

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
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
	}

	return cookie, nil
}
