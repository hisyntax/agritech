package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func (user *User) GetUserDoc(colName *mongo.Collection, filter interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data *User

	if err := colName.FindOne(ctx, filter).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
