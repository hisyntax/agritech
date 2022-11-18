package persistence

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// redis connection
var Cached = redis.NewClient(&redis.Options{
	Addr: os.Getenv("REDIS_ADDR"),
})

// set value in redis
func SetRedisValue(key string, value string, expiry time.Duration) error {
	errr := Cached.Set(key, value, expiry).Err()
	if errr != nil {
		return errr
	}
	return nil
}

func GetRedisValue(key string) (string, error) {
	value, err := Cached.Get(key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func DbInstance() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
	MongoDb := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("successfully connected to mongodb")
	return client
}

var Client *mongo.Client = DbInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	databaseName := os.Getenv("DATABASE_NAME")
	var collection *mongo.Collection = client.Database(databaseName).Collection(collectionName)
	return collection
}

func GetMongoDocs(colName *mongo.Collection, filter interface{}, opts ...*options.FindOptions) ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data []bson.M

	filterCusor, err := colName.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	if err := filterCusor.All(ctx, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func GetMongoDoc(colName *mongo.Collection, filter interface{}) (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data bson.M

	if err := colName.FindOne(ctx, filter).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func CreateMongoDoc(colName *mongo.Collection, data interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertNum, insertErr := colName.InsertOne(ctx, data)
	if insertErr != nil {
		return nil, insertErr
	}

	return insertNum, nil
}
func CreateManyMongoDoc(colName *mongo.Collection, data []interface{}) (*mongo.InsertManyResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertNum, insertErr := colName.InsertMany(ctx, data)
	if insertErr != nil {
		return nil, insertErr
	}

	return insertNum, nil
}

func UpdateMongoDoc(colName *mongo.Collection, filter interface{}, data interface{}) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateData := bson.D{{Key: "$set", Value: data}}
	res, err := colName.UpdateOne(ctx, filter, updateData)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func DeleteOneMongoDBDoc(colName *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := colName.DeleteOne(ctx, filter)

	if err != nil {
		return nil, err
	}
	ctx.Done()

	return res, nil
}

func DeleteManyMongoDBDoc(colName *mongo.Collection, filter []interface{}) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := colName.DeleteMany(ctx, filter)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func CountCollection(colName *mongo.Collection, filter interface{}) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := colName.CountDocuments(ctx, filter)

	if err != nil {
		return 0
	}
	return count
}
