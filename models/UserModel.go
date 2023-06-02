package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserPublic struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty""`
	FirstName string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName  string             `json:"last_name" bson:"last_name" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required"`
	Location  string             `json:"location" bson:"location" validate:"required"`
	RoleID    int                `json:"role_id,omitempty" bson:"role_id,omitempty"`
}

type NewUser struct {
	FirstName string `json:"first_name" bson:"first_name" validate:"required"`
	LastName  string `json:"last_name" bson:"last_name" validate:"required"`
	Email     string `json:"email" bson:"email" validate:"required"`
	Password  string `json:"password" bson:"password" validate:"required"`
	Location  string `json:"location" bson:"location" validate:"required"`
}

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty""`
	FirstName string             `json:"first_name" bson:"first_name" validate:"required"`
	LastName  string             `json:"last_name" bson:"last_name" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required"`
	Password  string             `json:"password" bson:"password" validate:"required"`
	Location  string             `json:"location" bson:"location" validate:"required"`
	RoleID    int                `json:"role_id,omitempty" bson:"role_id,omitempty"`
}
