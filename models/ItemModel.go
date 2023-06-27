package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

//TODO: Add Currency Enums

type Price struct {
	Amount   float64 `json:"amount" bson:"amount" validate:"required" example:"123"`
	Currency string  `json:"currency" bson:"currency"`
}

type ItemNew struct {
	Name        string `json:"name" bson:"name" validate:"required" example:"My Item"`
	Description string `json:"description" bson:"description" validate:"required" example:"My Item Description"`
	Price       Price  `json:"price" bson:"price" validate:"required"`
	Quantity    int    `json:"quantity" bson:"quantity" validate:"required" example:"123"`
}

type Item struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name" validate:"required" example:"My Item"`
	Description    string             `json:"description" bson:"description" validate:"required" example:"My Item Description"`
	Price          Price              `json:"price" bson:"price" validate:"required"`
	Quantity       int                `json:"quantity" bson:"quantity" validate:"required" example:"123"`
	OrganisationID primitive.ObjectID `json:"organisation_id" bson:"organisation_id" validate:"required"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
}
