package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/pkg/errors"
)

func Add(message string) int {
	collection := GetCollection("test")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", message}})
	if err != nil {
		return 0
	}
	fmt.Println("inserted on db")
	fmt.Println(res)
	return 0
}

func Save(element interface{},collection string) (string,error){
	col := GetCollection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := col.InsertOne(ctx, element)
	if err != nil{
		return "",errors.Wrap(err,"error saving element")
	}

	return res.InsertedID.(string),nil
}

func getClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		panic(err)
	}


	return client
}

func GetCollection(collectionName string) *mongo.Collection {
	client := getClient()

	collection := client.Database("testing").Collection(collectionName)

	return collection
}