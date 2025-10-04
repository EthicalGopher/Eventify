// Package db contains all the functions for mongodb
package db

import (
	"context"
	"eventify/datatype"
	"fmt"

	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var db *mongo.Collection
var client *mongo.Client

// Connect connects to mongodb
func Connect() {
	url := os.Getenv("MONGO_URL")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal("Could not ping MongoDb:", err)
	}
	db = client.Database("Eventify").Collection("user")
	fmt.Println("Mongo DB successfully Connected")
}

// Disconnect to gracefully disconnect from mongodb
func Disconnect() {
	err := client.Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}
}

// AddUsers it adds the user to the Database
func AddUsers(user datatype.User) (string, error) {
	result, err := db.InsertOne(context.TODO(), user)
	if err != nil {
		return "", fmt.Errorf("Error while InsertOne function : %v", err.Error())
	}
	return fmt.Sprintln("Inserted successfully at id : ", result.InsertedID), nil
}

// FindUser takes user id as parameter and finds the user
func FindUser(id string) (*datatype.User, error) {
	var user datatype.User
	err := db.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("Error while FindUser function : %v", err.Error())
	}
	return &user, nil
}
