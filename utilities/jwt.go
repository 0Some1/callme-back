package utilities

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const SecretKey = "salty"

func GenerateJwt(issuer string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // it will be last for 1 day
	})

	return claims.SignedString([]byte(SecretKey))

}

func ParseJwt(cookie string) (string, error) {

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}
	claims := token.Claims.(*jwt.StandardClaims)

	return claims.Issuer, nil

}
