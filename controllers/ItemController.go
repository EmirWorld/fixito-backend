package controllers

import (
	"context"
	"github.com/biter777/countries"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"poosible-backend/config"
	"poosible-backend/models"
	"poosible-backend/responses"
	"time"
)

var itemCollection *mongo.Collection = config.GetCollection(config.Database, "items")

// CreateItem godoc
// @Summary Create an item
// @Description Creates an item
// @Tags Item
// @Accept json
// @Produce json
// @Param item body models.ItemNew true "Item data"
// @Success 200 {object} responses.ItemResponse
// @Failure 400 {object} responses.ItemResponse
// @Failure 500 {object} responses.ItemResponse
// @Router /api/item [post]
// @Security BearerAuth
func CreateItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		var item *models.ItemNew
		var organisation *models.Organisation
		defer cancel()

		user, err := CurrentUser(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ItemResponse{Status: http.StatusInternalServerError, Message: "Error getting current user", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if user.OrganisationID == nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "User does not belong to an organisation", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		organisationId, _ := primitive.ObjectIDFromHex(user.OrganisationID.Hex())

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
		itemCheckErr := itemCollection.FindOne(ctx, bson.M{"name": item.Name, "organisation_id": organisationId}).Decode(&itemCheck)
		if itemCheckErr == nil {
			c.JSON(http.StatusBadRequest, responses.ItemResponse{Status: http.StatusBadRequest, Message: "Item name already exists", Data: map[string]interface{}{"data": itemCheck}})
			return
		}

		orgErr := organisationCollection.FindOne(ctx, bson.M{"_id": organisationId}).Decode(&organisation)
		if orgErr != nil {
			c.JSON(http.StatusInternalServerError, responses.OrganisationResponse{Status: http.StatusInternalServerError, Message: "Error getting organisation", Data: map[string]interface{}{"data": orgErr.Error()}})
			return
		}
		//get currency alpha
		currencies := countries.AllCurrenciesInfo()
		var currencyAlpha string
		for _, currency := range currencies {
			if int(currency.Code) == organisation.Currency {
				currencyAlpha = currency.Alpha
			}
		}

		price := models.Price{
			Amount:   item.Price.Amount,
			Currency: currencyAlpha,
		}

		newItem := models.Item{
			Name:           item.Name,
			Description:    item.Description,
			OrganisationID: organisationId,
			Price:          price,
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
		err := itemCollection.FindOne(ctx, bson.M{"_id": itemObjectId}).Decode(&item)
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
// @Success 200 {object} responses.ItemResponse
// @Failure 500 {object} responses.ItemResponse
// @Router /api/items [get]
// @Security BearerAuth
func GetItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		var items []*models.Item
		defer cancel()

		user, err := CurrentUser(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ItemResponse{Status: http.StatusInternalServerError, Message: "Internal server error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if user.OrganisationID == nil {
			c.JSON(http.StatusBadRequest, responses.ItemResponse{Status: http.StatusBadRequest, Message: "User does not belong to an organisation", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		organisationObjectId, _ := primitive.ObjectIDFromHex(user.OrganisationID.Hex())

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
