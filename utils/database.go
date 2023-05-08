package utils

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBInstance *mongo.Client = nil

func CheckBase() *mongo.Client {
	if MongoDBInstance == nil {
		MongoDBInstance = ConnectToBase()
		fmt.Printf("v funkciji")
	}
	return MongoDBInstance
}

// ConnectToBase creates a connection to the MongoDB instance and returns it
func ConnectToBase() *mongo.Client {
	newMongoInstance, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://"+url.QueryEscape("drvocepalci")+":"+url.QueryEscape("drvocepalci")+"@cluster0.ioplwkz.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		log.Printf("(ConnectToMongoDB) There was an error creating the mongoDB mongoDBInstance: %v", err)
	}

	log.Printf("(ConnectToMongoDB) Successfuly Connected to MongoDB")
	return newMongoInstance
}
