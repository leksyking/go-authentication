package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                          primitive.ObjectID `bson:"_id"`
	FirstName                   *string            `json:"firstname" validate:"required"`
	LastName                    *string            `json:"lastname" validate:"required"`
	Email                       *string            `json:"email" validate:"required,email"`
	UserName                    *string            `json:"username" validate:"required"`
	Password                    *string            `json:"password" validate:"required"`
	VerificationToken           *string            `json:"verification_token" bson:"verification_token"`
	IsVerified                  *bool              `json:"is_verified" bson:"is_verified"`
	Verified                    time.Time          `json:"verified" bson:"verified"`
	PasswordToken               string             `json:"passwordtoken"`
	PasswordTokenExpirationDate int64              `json:"passwordtokenexpirationdate"`
}

type Token struct {
	TokenID      primitive.ObjectID `bson:"_id"`
	RefreshToken *string            `json:"refreshtoken" bson:"refreshtoken"`
	IP           *string            `json:"ip" bson:"ip"`
	IsValid      bool               `json:"is_valid" bson:"is_valid"`
	UserAgent    *string            `json:"useragent" bson:"useragent"`
	User         primitive.ObjectID `json:"user_id" bson:"user_id"`
}
