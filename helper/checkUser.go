package helper

import (
	"github.com/hisyntax/agritech/persistence"

	"go.mongodb.org/mongo-driver/bson"
)

func CheckUser(email string) (string, error) {
	userFilter := bson.D{{Key: "email", Value: email}}
	_, err := persistence.GetMongoDoc(persistence.UserCollection, userFilter)
	if err != nil {
		return "", err
	}

	return email, nil
}
