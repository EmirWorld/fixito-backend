package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authSession := sessions.Default(c)

		// Check if user is authenticated
		isAuthenticated := authSession.Get("authenticated")
		if isAuthenticated != true {
			c.AbortWithStatusJSON(401, gin.H{"message": "User is not logged in"})
			return
		}

		// Check if access token is valid
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"message": "Missing Authorization Header"})
			return
		}

		//Remove the Bearer prefix from the token
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		//Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, nil
			}
			return []byte("secret_key"), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid Token"})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"message": "Invalid Token"})
			return
		}

		// Check if token is for current user
		userIDFromToken := token.Claims.(jwt.MapClaims)["user_id"].(string)
		currentUserID := authSession.Get("user_id").(string)

		if userIDFromToken != currentUserID {
			c.AbortWithStatusJSON(401, gin.H{"message": "Token is not for current user"})
			return
		}

		c.Next()
	}
}
