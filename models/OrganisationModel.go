package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Organisation struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name" validate:"required" example:"My Organization"`
}
