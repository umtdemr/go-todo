package user

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateNewJWT(t *testing.T) {
	username := "username"
	tokenString, err := GenerateNewJWT(username)

	assert.Nil(t, err)

	token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected siging method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	assert.Nil(t, parseErr)

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		_, isExpirationExists := claims["exp"] // no need to check if the token expired since it's automatically checked
		parsedUsername, _ := claims["username"]

		assert.True(t, isExpirationExists)
		assert.Equal(t, username, parsedUsername)
	}
}

func TestValidateJWT(t *testing.T) {
	validToken, _ := GenerateNewJWT("testuser")

	tests := []struct {
		name          string
		tokenString   string
		expectedError error
		expectedUser  string
	}{
		{
			name:          "Valid Token",
			tokenString:   validToken,
			expectedError: nil,
			expectedUser:  "testuser",
		},
		{
			name:          "Invalid Token",
			tokenString:   "invalid.token.string",
			expectedError: ErrTokenNotValid,
			expectedUser:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user, err := ValidateJWT(tc.tokenString)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedUser, user)
		})
	}
}
