package user

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtSecret = []byte("TOP_SECRET_KEY_FOR_JWT")

func GenerateNewJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)

	return tokenString, err
}

func ValidateJWT(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected siging method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return false, err
	}
	if !token.Valid {
		return false, ErrorTokenNotValid
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		_, isExpirationExists := claims["exp"] // no need to check if the token expired since it's automatically checked
		if !isExpirationExists {
			return false, ErrorTokenNotValid
		}

		_, isUsernameExistsOnMap := claims["username"]
		if !isUsernameExistsOnMap {
			return false, ErrorTokenNotValid
		}

		return true, nil
	}
	return false, ErrorTokenNotValid
}
