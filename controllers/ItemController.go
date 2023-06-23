package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"possible-backend/config"
	"possible-backend/models"
	"possible-backend/responses"
	"time"
)

var itemCollection *mongo.Collection = config.GetCollection(config.Database, "items")

// CreateItem godoc
// @Summary Create an item
// @Description Creates an item
// @Tags Item
// @Accept json
// @Produce json
// @Param organisationId path string true "Organisation ID"
// @Param item body models.ItemNew true "Item data"
// @Success 200 {object} responses.ItemResponse
// @Failure 400 {object} responses.ItemResponse
// @Failure 500 {object} responses.ItemResponse
// @Router /api/item/{organisationId} [post]
// @Security BearerAuth
func CreateItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		var item *models.ItemNew
		organisationId := c.Param("organisationId")
		organisationObjectId, _ := primitive.ObjectIDFromHex(organisationId)
		defer cancel()

		// Bind the JSON body to the struct
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, responses.ItemResponse{Status: http.StatusBadRequest, Message: "Bad request", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Validate the input
		if err := validate.Struct(item); err != nil {
			c.JSON(http.StatusBadRequest, responses.ItemResponse{Status: http.StatusBadRequest, Message: "Bad request", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//Check if item name exists
		var itemCheck models.Item
		err := itemCollection.FindOne(ctx, models.Item{Name: item.Name}).Decode(&itemCheck)
		if err == nil {
			c.JSON(http.StatusBadRequest, responses.ItemResponse{Status: http.StatusBadRequest, Message: "Item name already exists", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		newItem := models.Item{
			Name:           item.Name,
			Description:    item.Description,
			OrganisationID: organisationObjectId,
			Price:          item.Price,
			Quantity:       item.Quantity,
			UpdatedAt:      time.Now(),
			CreatedAt:      time.Now(),
		}

		// Insert the organisation into the database
		result, err := itemCollection.InsertOne(ctx, newItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ItemResponse{Status: http.StatusInternalServerError, Message: "Internal server error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Return the organisation
		c.JSON(http.StatusOK, responses.ItemResponse{Status: http.StatusOK, Message: "Item created", Data: map[string]interface{}{"data": result}})

	}
}

// GetItem godoc
// @Summary Get an item
// @Description Gets an item
// @Tags Item
// @Accept json
// @Produce json
// @Param itemId path string true "Item ID"
// @Success 200 {object} responses.ItemResponse
// @Failure 404 {object} responses.ItemResponse
// @Failure 500 {object} responses.ItemResponse
// @Router /api/item/{itemId} [get]
// @Security BearerAuth
func GetItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		var item *models.Item
		itemId := c.Param("itemId")
		itemObjectId, _ := primitive.ObjectIDFromHex(itemId)
		defer cancel()

		// Find the item in the database
		err := itemCollection.FindOne(ctx, models.Item{ID: itemObjectId}).Decode(&item)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, responses.ItemResponse{Status: http.StatusNotFound, Message: "Item not found", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			c.JSON(http.StatusInternalServerError, responses.ItemResponse{Status: http.StatusInternalServerError, Message: "Internal server error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Return the item
		c.JSON(http.StatusOK, responses.ItemResponse{Status: http.StatusOK, Message: "Item found", Data: map[string]interface{}{"data": item}})
	}
}

// GetItems godoc
// @Summary Get items
// @Description Gets items
// @Tags Item
// @Accept json
// @Produce json
// @Param organisationId path string true "Organisation ID"
// @Success 200 {object} responses.ItemResponse
// @Failure 500 {object} responses.ItemResponse
// @Router /api/items/{organisationId} [get]
// @Security BearerAuth
func GetItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		var items []*models.Item
		organisationId := c.Param("organisationId")
		organisationObjectId, _ := primitive.ObjectIDFromHex(organisationId)
		defer cancel()

		// Find the items in the database
		result, err := itemCollection.Find(ctx, bson.M{"organisation_id": organisationObjectId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ItemResponse{Status: http.StatusInternalServerError, Message: "Internal server error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Iterate through the result and read in optimal chunks
		defer result.Close(ctx)
		for result.Next(ctx) {
			var item *models.Item
			err := result.Decode(&item)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.ItemResponse{Status: http.StatusInternalServerError, Message: "Internal server error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
			items = append(items, item)
		}
		// Return the items
		c.JSON(http.StatusOK, responses.ItemResponse{Status: http.StatusOK, Message: "Items found", Data: map[string]interface{}{"data": items}})
	}
}

func UpdateItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func DeleteItem() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
