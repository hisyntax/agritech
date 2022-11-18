package user

import (
	"time"

	"github.com/hisyntax/agritech/helper"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Full_Name  string             `json:"full_name" bson:"full_name" validate:"required"`
	Email      string             `json:"email" bson:"email" validate:"required"`
	Password   string             `json:"password" bson:"password" validate:"required"`
	Created_At time.Time          `json:"created_at" bson:"created_at"`
	Updated_At time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserSignin struct {
	Email    string `json:"email" bson:"email" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required"`
}

type PublicUser struct {
	Full_Name string `json:"full_name" bson:"full_name" validate:"required"`
	Email     string `json:"email" bson:"email" validate:"required"`
}

func (User) PublicUser(user *User) *PublicUser {
	return &PublicUser{
		Full_Name: user.Full_Name,
		Email:     user.Email,
	}
}

func (u *User) BeforeCreate() error {
	hashPassword, err := helper.Hash(u.Password)
	if err != nil {
		return err
	}

	u.Password = string(hashPassword)
	return nil
}
