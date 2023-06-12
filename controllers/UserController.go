package controllers

import (
	"context"
	"fixito-backend/config"
	"fixito-backend/models"
	"fixito-backend/responses"
	"fixito-backend/validators"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var userCollection *mongo.Collection = config.GetCollection(config.Database, "users")
var validate = validator.New()

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.NewUser true "Create User"
// @Success 200 {object} responses.UserResponse
// @Failure 400 {object} responses.UserResponse
// @Failure 500 {object} responses.UserResponse
// @Router /api/user [post]
func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user *models.NewUser
		defer cancel()

		//Validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid request body", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid input fields", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//Validate password length and complexity
		errorMessage, ok := validators.ValidatePassword(user.Password)
		if !ok {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid password", Data: map[string]interface{}{"data": errorMessage}})
			return
		}

		//Check if email already exists
		var existingUser *models.UserPublic
		err := userCollection.FindOne(ctx, models.UserPublic{Email: user.Email}).Decode(&existingUser)
		if err != nil {
			if err != mongo.ErrNoDocuments {
				c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error checking if email already exists", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Email already exists", Data: map[string]interface{}{"data": nil}})
			return
		}

		//Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error hashing password", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//Insert user into database
		newUser := models.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  string(hashedPassword),
			Location:  user.Location,
			RoleID:    models.USER,
		}

		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error inserting user into database", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "User created successfully", Data: map[string]interface{}{"data": result}})
	}
}

// GetUser retrieves a user by ID.
// @Summary Get user
// @Description Retrieves a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} responses.UserResponse
// @Failure 400 {object} responses.UserResponse
// @Failure 404 {object} responses.UserResponse
// @Failure 500 {object} responses.UserResponse
// @Router /api/user/{userId} [get]
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		var user *models.UserPublic
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(userId)
		fmt.Println(objId)
		fmt.Println(userId)

		err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNotFound, Message: "User not found", Data: map[string]interface{}{"data": nil}})
				return
			}
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error getting user", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "User retrieved successfully", Data: map[string]interface{}{"data": user}})
	}
}
