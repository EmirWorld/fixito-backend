package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"possible-backend/config"
	"possible-backend/helpers"
	"possible-backend/models"
	"possible-backend/responses"
	"possible-backend/validators"
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
// @Param user body models.UserNew true "Create User"
// @Success 200 {object} responses.UserResponse
// @Failure 400 {object} responses.UserResponse
// @Failure 500 {object} responses.UserResponse
// @Router /api/user [post]
func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user *models.UserNew
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
// @Security BearerAuth
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

// UpdateUser updates a user with the given user ID.
// @Summary Update a user
// @Description Update a user with the given user ID
// @Tags Users
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param user body models.UserNew true "User object to update"
// @Success 200 {object} responses.UserResponse
// @Failure 400 {object} responses.UserResponse
// @Failure 500 {object} responses.UserResponse
// @Router /api/user/{userId} [put]
// @Security BearerAuth
func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		var user *models.UserNew
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(userId)

		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid input fields", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//Validate password length and complexity
		errorMessage, ok := validators.ValidatePassword(user.Password)
		if !ok {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid password", Data: map[string]interface{}{"data": errorMessage}})
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error hashing password", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//Update user in database
		update := bson.M{"firstName": user.FirstName, "lastName": user.LastName, "email": user.Email, "password": string(hashedPassword), "location": user.Location}
		result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error updating user", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedUser *models.UserPublic
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error getting updated user", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "User updated successfully", Data: map[string]interface{}{"data": updatedUser}})
	}
}

// DeleteUser deletes a user with the given user ID.
// @Summary Delete a user
// @Description Delete a user with the given user ID
// @Tags Users
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} responses.UserResponse
// @Failure 404 {object} responses.UserResponse
// @Failure 500 {object} responses.UserResponse
// @Router /api/user/{userId} [delete]
// @Security BearerAuth
func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(userId)

		//TODO: Delete user's organisation and items

		result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error deleting user", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNotFound, Message: "User not found", Data: map[string]interface{}{"data": nil}})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "User deleted successfully", Data: map[string]interface{}{"data": nil}})
	}
}

// GetCurrentUser
// @Summary Get current user
// @Description Get the details of the session user
// @Tags Users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} responses.UserResponse
// @Failure 400 {object} responses.UserResponse
// @Failure 401 {object} responses.UserResponse
// @Failure 500 {object} responses.UserResponse
// @Router /api/user/current [get]
// @Security BearerAuth
func GetCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := helpers.GetClaims(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid token", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userID, ok := claims["user_id"].(string)
		objId, _ := primitive.ObjectIDFromHex(userID)
		defer cancel()
		if !ok {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid token", Data: map[string]interface{}{"data": "Invalid  user claims"}})
			return
		}
		var user *models.UserPublic

		err = userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error getting user", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}})
	}
}
