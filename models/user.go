package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName *string            `json:"firstname" validator:"required"`
	LastName  *string            `json:"lastname" validator:"required"`
	Email     *string            `json:"email" validator:"required,email"`
	UserName  *string            `json:"username" validator:"required"`
	Password  *string            `json:"password" validator:"required"`
}
