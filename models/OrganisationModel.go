package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OrganisationNew struct {
	Name        string `json:"name" bson:"name" validate:"required" example:"My Organization"`
	Description string `json:"description" bson:"description" validate:"required" example:"My Organization Description"`
	Phone       string `json:"phone" bson:"phone" validate:"required" example:"1234567890"`
	Address     string `json:"address" bson:"address" validate:"required" example:"123 Main St"`
	Logo        string `json:"logo" bson:"logo" validate:"required" example:"https://www.example.com/logo.png"`
}

type Organisation struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" validate:"required" example:"My Organization"`
	Description string             `json:"description" bson:"description" validate:"required" example:"My Organization Description"`
	Phone       string             `json:"phone" bson:"phone" validate:"required" example:"1234567890"`
	Address     string             `json:"address" bson:"address" validate:"required" example:"123 Main St"`
	Logo        string             `json:"logo" bson:"logo" validate:"required" example:"https://www.example.com/logo.png"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}
