package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OrganisationNew struct {
	Name        string `json:"name" bson:"name" validate:"required" example:"My Organization"`
	Description string `json:"description" bson:"description" validate:"required" example:"My Organization Description"`
	Phone       string `json:"phone" bson:"phone" validate:"required" example:"1234567890"`
	Country     int    `json:"country_code" bson:"country_code" validate:"required" example:"36"`
	Address     string `json:"address" bson:"address" validate:"required" example:"123 Main St"`
	ZipCode     string `json:"zip_code" bson:"zip_code" validate:"required" example:"12345"`
	Logo        string `json:"logo" bson:"logo"  example:"https://www.example.com/logo.png"`
	Currency    int    `json:"currency_code" bson:"currency_code" validate:"required" example:"840"`
}

type OrganisationPublic struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" validate:"required" example:"My Organization"`
	Description string             `json:"description" bson:"description" validate:"required" example:"My Organization Description"`
	Phone       string             `json:"phone" bson:"phone" validate:"required" example:"1234567890"`
	Address     string             `json:"address" bson:"address" validate:"required" example:"123 Main St"`
	Country     string             `json:"country" bson:"country" validate:"required" example:"US"`
	Currency    string             `json:"currency" bson:"currency" validate:"required" example:"USD"`
	ZipCode     string             `json:"zip_code" bson:"zip_code" validate:"required" example:"12345"`
	Logo        string             `json:"logo" bson:"logo" example:"https://www.example.com/logo.png"`
}

type Organisation struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" validate:"required" example:"My Organization"`
	Description string             `json:"description" bson:"description" validate:"required" example:"My Organization Description"`
	Phone       string             `json:"phone" bson:"phone" validate:"required" example:"1234567890"`
	Address     string             `json:"address" bson:"address" validate:"required" example:"123 Main St"`
	Country     int                `json:"country_code" bson:"country_code" validate:"required" example:"36"`
	Currency    int                `json:"currency_code" bson:"currency_code" validate:"required" example:"840"`
	ZipCode     string             `json:"zip_code" bson:"zip_code" validate:"required" example:"12345"`
	Logo        string             `json:"logo" bson:"logo" validate:"required" example:"https://www.example.com/logo.png"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}
