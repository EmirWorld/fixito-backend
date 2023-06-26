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
	"poosible-backend/helpers"
	"poosible-backend/models"
	"poosible-backend/responses"
	"time"
)

var organisationCollection *mongo.Collection = config.GetCollection(config.Database, "organisations")

// CreateOrganisation godoc
// @Summary Create an organisation
// @Description Creates an organisation
// @Tags Organisation
// @Accept json
// @Produce json
// @Param organisation body models.OrganisationNew true "Organisation data"
// @Success 200 {object} responses.OrganisationResponse
// @Failure 400 {object} responses.OrganisationResponse
// @Failure 500 {object} responses.OrganisationResponse
// @Router /api/organisation [post]
// @Security BearerAuth
func CreateOrganisation() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := helpers.GetClaims(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.OrganisationResponse{Status: http.StatusBadRequest, Message: "Bad request", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		userID, ok := claims["user_id"].(string)
		userObjectID, _ := primitive.ObjectIDFromHex(userID)
		var organisation *models.OrganisationNew
		defer cancel()

		if !ok {
			c.JSON(http.StatusBadRequest, responses.OrganisationResponse{Status: http.StatusBadRequest, Message: "Bad request", Data: map[string]interface{}{"data": "Invalid user id"}})
			return
		}

		// Bind the JSON body to the struct
		if err := c.ShouldBindJSON(&organisation); err != nil {
			c.JSON(http.StatusBadRequest, responses.OrganisationResponse{Status: http.StatusBadRequest, Message: "Bad request", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Validate the input
		if err := validate.Struct(organisation); err != nil {
			c.JSON(http.StatusBadRequest, responses.OrganisationResponse{Status: http.StatusBadRequest, Message: "Bad request", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//check if organisation name is empty
		if organisation.Name == "" {
			c.JSON(http.StatusBadRequest, responses.OrganisationResponse{Status: http.StatusBadRequest, Message: "Bad request", Data: map[string]interface{}{"data": "Organisation name is required"}})
			return
		}

		//check if user already has an organisation
		var user *models.User
		err = userCollection.FindOne(ctx, bson.M{"_id": userObjectID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.OrganisationResponse{Status: http.StatusInternalServerError, Message: "Error getting user", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if user.OrganisationID != nil {
			c.JSON(http.StatusBadRequest, responses.OrganisationResponse{Status: http.StatusBadRequest, Message: "Bad request", Data: map[string]interface{}{"data": "User already has an organisation"}})
			return
		}

		newOrganisation := models.Organisation{
			Name:        organisation.Name,
			Description: organisation.Description,
			Phone:       organisation.Phone,
			Address:     organisation.Address,
			Logo:        organisation.Logo,
			Country:     organisation.Country,
			Currency:    organisation.Currency,
			ZipCode:     organisation.ZipCode,
			UpdatedAt:   time.Now(),
			CreatedAt:   time.Now(),
		}

		result, err := organisationCollection.InsertOne(ctx, newOrganisation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.OrganisationResponse{Status: http.StatusInternalServerError, Message: "Error creating organisation", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.OrganisationResponse{Status: http.StatusOK, Message: "Organisation created successfully", Data: map[string]interface{}{"data": result}})

		//update user organization ID in database
		update := bson.M{"$set": bson.M{"organisation_id": result.InsertedID}}
		_, err = userCollection.UpdateOne(ctx, bson.M{"_id": userObjectID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.OrganisationResponse{Status: http.StatusInternalServerError, Message: "Error settings organization for user", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

	}
}

// GetOrganisation godoc
// @Summary Get an organisation
// @Description Gets an organisation
// @Tags Organisation
// @Accept json
// @Produce json
// @Param organisationId path string true "Organisation ID"
// @Success 200 {object} responses.OrganisationResponse
// @Failure 400 {object} responses.OrganisationResponse
// @Failure 500 {object} responses.OrganisationResponse
// @Router /api/organisation/{organisationId} [get]
// @Security BearerAuth
func GetOrganisation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		organisationID := c.Param("organisationId")
		organisationObjectID, _ := primitive.ObjectIDFromHex(organisationID)
		var organisation *models.Organisation
		defer cancel()

		err := organisationCollection.FindOne(ctx, bson.M{"_id": organisationObjectID}).Decode(&organisation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.OrganisationResponse{Status: http.StatusInternalServerError, Message: "Error getting organisation", Data: map[string]interface{}{"data": err.Error()}})
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

		organisationPublic := models.OrganisationPublic{
			ID:          organisation.ID,
			Name:        organisation.Name,
			Description: organisation.Description,
			Phone:       organisation.Phone,
			Logo:        organisation.Logo,
			Country:     countries.ByNumeric(organisation.Country).Info().Name,
			Address:     organisation.Address,
			ZipCode:     organisation.ZipCode,
			Currency:    currencyAlpha,
		}

		c.JSON(http.StatusOK, responses.OrganisationResponse{Status: http.StatusOK, Message: "Organisation retrieved successfully", Data: map[string]interface{}{"data": organisationPublic}})
	}
}

// UpdateOrganisation godoc
// @Summary Update an organisation
// @Description Updates an organisation
// @Tags Organisation
// @Accept json
// @Produce json
// @Param organisationId path string true "Organisation ID"
// @Param organisation body models.OrganisationNew true "Organisation object to be updated"
// @Success 200 {object} responses.OrganisationResponse
// @Failure 400 {object} responses.OrganisationResponse
// @Failure 500 {object} responses.OrganisationResponse
// @Router /api/organisation/{organisationId} [put]
// @Security BearerAuth
func UpdateOrganisation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		organisationID := c.Param("organisationId")
		organisationObjectID, _ := primitive.ObjectIDFromHex(organisationID)
		var organisation *models.OrganisationNew
		defer cancel()

		// Bind the JSON body to the struct
		if err := c.ShouldBindJSON(&organisation); err != nil {
			c.JSON(http.StatusBadRequest, responses.OrganisationResponse{Status: http.StatusBadRequest, Message: "Bad request", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// Validate the input
		if err := validate.Struct(organisation); err != nil {
			c.JSON(http.StatusBadRequest, responses.OrganisationResponse{Status: http.StatusBadRequest, Message: "Bad request", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		update := bson.M{"$set": bson.M{"name": organisation.Name, "country_code": organisation.Country, "currency_code": organisation.Country, "zip_code": organisation.ZipCode, "description": organisation.Description, "logo": organisation.Logo, "address": organisation.Address, "phone": organisation.Phone, "updated_at": time.Now()}}
		_, err := organisationCollection.UpdateOne(ctx, bson.M{"_id": organisationObjectID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.OrganisationResponse{Status: http.StatusInternalServerError, Message: "Error updating organisation", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated organisation
		c.JSON(http.StatusOK, responses.OrganisationResponse{Status: http.StatusOK, Message: "Organisation updated successfully", Data: map[string]interface{}{"data": organisation}})

	}
}

// DeleteOrganisation godoc
// @Summary Delete an organisation
// @Description Deletes an organisation
// @Tags Organisation
// @Accept json
// @Produce json
// @Param organisationId path string true "Organisation ID"
// @Success 200 {object} responses.OrganisationResponse
// @Failure 400 {object} responses.OrganisationResponse
// @Failure 500 {object} responses.OrganisationResponse
// @Router /api/organisation/{organisationId} [delete]
// @Security BearerAuth
func DeleteOrganisation() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get organisation ID
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		organisationID := c.Param("organisationId")
		organisationObjectID, _ := primitive.ObjectIDFromHex(organisationID)
		defer cancel()

		//TODO: check if organisation has users and delete id from users

		//Delete organisation
		_, err := organisationCollection.DeleteOne(ctx, bson.M{"_id": organisationObjectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.OrganisationResponse{Status: http.StatusInternalServerError, Message: "Error deleting organisation", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.OrganisationResponse{Status: http.StatusOK, Message: "Organisation deleted successfully", Data: map[string]interface{}{"data": nil}})
	}
}
