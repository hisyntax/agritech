package persistence

import (
	"errors"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	UserCollection *mongo.Collection = OpenCollection(Client, "users")
	OtpCollection  *mongo.Collection = OpenCollection(Client, "otp")

	Validate = validator.New()
)

var (
	ErrEmailTaken        = errors.New("email taken")
	ErrEmailNotFound     = errors.New("email not found")
	ErrPasswordIncorrect = errors.New("password incorrect")
)
