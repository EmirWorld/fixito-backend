package helpers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"possible-backend/responses"
	"strings"
)

func GetClaims(c *gin.Context) (jwt.MapClaims, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid token", Data: map[string]interface{}{"data": "Invalid token"}})
		return nil, nil
	}

	//remove the Bearer prefix from the token
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte("secret_key"), nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid token", Data: map[string]interface{}{"data": err.Error()}})
		return nil, err
	}

	if !token.Valid {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid token", Data: map[string]interface{}{"data": "Invalid token"}})
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid token", Data: map[string]interface{}{"data": "Invalid token claims"}})
		return nil, err
	}

	return claims, nil
}
