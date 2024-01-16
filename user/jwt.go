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

func GenerateResetPasswordToken(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":                  email,
		"exp":                    expirationTime.Unix(),
		"isRequestPasswordToken": true, // identifier. We seek if the props exists. Value is not important here
	})

	tokenString, err := token.SignedString(jwtSecret)

	return tokenString, err
}

func ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected siging method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", ErrTokenNotValid
	}
	if !token.Valid {
		return "", ErrTokenNotValid
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		_, isExpirationExists := claims["exp"] // no need to check if the token expired since it's automatically checked
		if !isExpirationExists {
			return "", ErrTokenNotValid
		}

		username, isUsernameExistsOnMap := claims["username"]
		if !isUsernameExistsOnMap {
			return "", ErrTokenNotValid
		}

		return username.(string), nil
	}
	return "", ErrTokenNotValid
}

// TODO: refactor validating JWT
func ValidateResetPasswordToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected siging method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", ErrTokenNotValid
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		_, isExpirationExists := claims["exp"] // no need to check if the token expired since it's automatically checked
		if !isExpirationExists {
			return "", ErrTokenNotValid
		}

		_, isRequestPasswordTokenErr := claims["isRequestPasswordToken"]

		if !isRequestPasswordTokenErr {
			return "", ErrTokenNotValid
		}

		email, isEmailExistsOnMap := claims["email"]
		if !isEmailExistsOnMap {
			return "", ErrTokenNotValid
		}

		return email.(string), nil
	}
	return "", ErrTokenNotValid
}
