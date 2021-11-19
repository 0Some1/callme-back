package lib

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const SecretKey = "salty"

func GenerateJwt(issuer string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(time.Hour * 99999).Unix(), // it will be last for 1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		fmt.Print("lib.GenerateJwt - ")
		return "", err
	}

	return token, nil

}

func ParseJwt(cookie string) (string, error) {

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		fmt.Print("lib.ParseJwt - ")
		return "", err
	}
	claims := token.Claims.(*jwt.StandardClaims)

	return claims.Issuer, nil

}
