// Package db contains all the functions for mongodb
package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var db *mongo.Collection

type User struct {
	ID        string   `bson:"_id,omitempty"`
	Name      string   `bson:"name"`
	Email     string   `bson:"email"`
	Password  string   `bson:"password"`
	Interests []string `bson:"interests,omitempty"`
	Role      string   `bson:"role"`
}

// Connect connects to mongodb
func Connect() {
	url := os.Getenv("MONGO_URL")
	fmt.Println(url)
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

// AddUsers it adds the user to the Database
func AddUsers(user User) {
	result, err := db.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal("Error while InsertOne function : ", err)
	}
	fmt.Println("Inserted successfully at id : ", result.InsertedID)
}
