package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("secretpassword")

func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // token is valid from 1 hour

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // HS256
	return token.SignedString(secretKey)
}

func Verifytoken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token")
}
