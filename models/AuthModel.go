package models

type Credentials struct {
	Email    string `json:"email" bson:"email" validate:"required" example:"emirkovacevic@protonmail.com"`
	Password string `json:"password" bson:"password" validate:"required" example:"408660As!"`
}
