package controllers

import (
	"context"
	"github.com/biter777/countries"
	"github.com/gin-gonic/gin"
	"net/http"
	"poosible-backend/responses"
	"time"
)

// Currencies godoc
// @Summary Get all currencies
// @Description Get all currencies
// @Tags Helper
// @Accept json
// @Produce json
// @Success 200 {object} responses.ItemResponse
// @Failure 400 {object} responses.ItemResponse
// @Failure 500 {object} responses.ItemResponse
// @Router /api/helper/currencies [get]
func Currencies() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		currencies := countries.AllCurrenciesInfo()
		defer cancel()

		if len(currencies) == 0 {
			c.JSON(http.StatusInternalServerError, responses.HelperResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": "No currencies found"}})
			return
		}

		c.JSON(http.StatusOK, responses.HelperResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": currencies, "count": len(currencies)}})
	}
}

// Countries godoc
// @Summary Get all countries
// @Description Get all countries
// @Tags Helper
// @Accept json
// @Produce json
// @Success 200 {object} responses.ItemResponse
// @Failure 400 {object} responses.ItemResponse
// @Failure 500 {object} responses.ItemResponse
// @Router /api/helper/countries [get]
func Countries() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		countries := countries.AllCapitalsInfo()
		defer cancel()

		if len(countries) == 0 {
			c.JSON(http.StatusInternalServerError, responses.HelperResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": "No countries found"}})
			return
		}

		c.JSON(http.StatusOK, responses.HelperResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": countries, "count": len(countries)}})
	}
}
