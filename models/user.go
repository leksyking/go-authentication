package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id"`
	FirstName         *string            `json:"firstname" validate:"required"`
	LastName          *string            `json:"lastname" validate:"required"`
	Email             *string            `json:"email" validate:"required,email"`
	UserName          *string            `json:"username" validate:"required"`
	Password          *string            `json:"password" validate:"required"`
	VerificationToken *string            `json:"verification_token"`
	IsVerified        *bool              `json:"is_verified"`
	Verified          time.Time          `json:"verified"`
}

type Token struct {
	TokenID      primitive.ObjectID `bson:"_id"`
	RefreshToken *string            `json:"refreshtoken"`
	IP           *string            `json:"ip"`
	IsValid      bool               `json:"is_valid"` //set default to true
	UserAgent    *string            `json:"useragent"`
	User         primitive.ObjectID `json:"user_id"`
}
