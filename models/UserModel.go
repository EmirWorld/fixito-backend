package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserPublic struct {
	ID             primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName      string              `json:"first_name" bson:"first_name" validate:"required"`
	LastName       string              `json:"last_name" bson:"last_name" validate:"required"`
	Email          string              `json:"email" bson:"email" validate:"required"`
	Location       string              `json:"location" bson:"location" validate:"required"`
	OrganisationID *primitive.ObjectID `json:"organization_id,omitempty" bson:"organization_id,omitempty"`
	RoleID         int                 `json:"role_id,omitempty" bson:"role_id,omitempty"`
}

type NewUser struct {
	FirstName string `json:"first_name" bson:"first_name" validate:"required" example:"John"`
	LastName  string `json:"last_name" bson:"last_name" validate:"required" example:"Doe"`
	Email     string `json:"email" bson:"email" validate:"required" example:"johndoe@mail.com"`
	Password  string `json:"password" bson:"password" validate:"required" example:"password123"`
	Location  string `json:"location" bson:"location" validate:"required" example:"New York"`
}

type User struct {
	ID             primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName      string              `json:"first_name" bson:"first_name" validate:"required"`
	LastName       string              `json:"last_name" bson:"last_name" validate:"required"`
	Email          string              `json:"email" bson:"email" validate:"required"`
	Password       string              `json:"password" bson:"password" validate:"required"`
	Location       string              `json:"location" bson:"location" validate:"required"`
	OrganisationID *primitive.ObjectID `json:"organisation_id,omitempty" bson:"organisation_id,omitempty"`
	RoleID         int                 `json:"role_id,omitempty" bson:"role_id,omitempty"`
}
