package controllers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	tokenString, err := token.SignedString([]byte("your_secret_key")) // Replace with your secret key
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
