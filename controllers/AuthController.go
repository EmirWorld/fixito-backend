package controllers

import (
	"context"
	"fixito-backend/models"
	"fixito-backend/responses"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// generateToken generates a JWT token and stores it in the session
func generateToken(user models.User, c *gin.Context) error {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = user.ID.Hex()
	expirationTime := time.Now().Add(time.Hour * 24) // Token will be valid for 24 hours
	claims["exp"] = expirationTime.Unix()
	refreshUUID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	refreshUUIDString := refreshUUID.String()
	claims["refresh_uuid"] = refreshUUIDString
	tokenString, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		return err
	}

	authSession := sessions.Default(c)
	authSession.Set("authenticated", true)
	authSession.Set("access_token", tokenString)
	authSession.Set("access_token_expiry", expirationTime)
	authSession.Set("refresh_token", refreshUUIDString)

	// Save it before we write to the response/return from the handler.
	err = authSession.Save()
	if err != nil {
		return err
	}
	return nil
}

// Login godoc
// @Summary Log in a user
// @Description Authenticates a user and returns a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.Credentials true "User credentials"
// @Success 200 {object} responses.AuthResponse
// @Failure 400 {object} responses.AuthResponse
// @Failure 401 {object} responses.AuthResponse
// @Failure 500 {object} responses.AuthResponse
// @Router /api/auth/login [post]
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials *models.Credentials

		//Validate the request body
		if err := c.BindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid request body", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//Use validators library to validate the email and password
		if err := validate.Struct(credentials); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid input fields", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var user *models.User
		err := userCollection.FindOne(context.Background(), bson.M{"email": credentials.Email}).Decode(&user)
		fmt.Println(user.Password)
		fmt.Println(credentials.Password)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid email or password", Data: map[string]interface{}{"data": nil}})
				return
			}
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error checking if email already exists", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//Compare the stored hashed password, with the hashed version of the password that was received
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid email or password", Data: map[string]interface{}{"data": nil}})
			return
		}

		//Generate JWT token
		err = generateToken(*user, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error generating token", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//Get token from session
		authSession := sessions.Default(c)

		token := authSession.Get("access_token")
		expiresAt := authSession.Get("access_token_expiry")
		refreshToken := authSession.Get("refresh_token")
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "Login successful", Data: map[string]interface{}{"data": map[string]interface{}{"token": token, "expires_at": expiresAt, "refresh_token": refreshToken}}})

	}
}

// Logout godoc
// @Summary Logout a user
// @Description Logs out a user and clears the session
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} responses.AuthResponse
// @Failure 400 {object} responses.AuthResponse
// @Failure 401 {object} responses.AuthResponse
// @Failure 500 {object} responses.AuthResponse
// @Router /api/auth/logout [post]
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		//delete session
		authSessions := sessions.Default(c)
		authSessions.Clear()
		err := authSessions.Save()
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.AuthResponse{Status: http.StatusInternalServerError, Message: "Error deleting session", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.AuthResponse{Status: http.StatusOK, Message: "Logout successful", Data: nil})
	}
}
