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

var client *mongo.Client

// Connect connects to mongodb
func Connect() {
	url := os.Getenv("MONGO_URL")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)
	var err error
	client, err = mongo.Connect(opts)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal("Could not ping MongoDb:", err)
	}
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
	db := client.Database("eventify").Collection("user")
	result, err := db.InsertOne(context.TODO(), user)
	if err != nil {
		return "", fmt.Errorf("Error while AddUsers function : %v", err.Error())
	}
	return fmt.Sprintln("Inserted successfully at id : ", result.InsertedID), nil
}

// FindUser takes user id as parameter and finds the user
func FindUser(id string) (*datatype.User, error) {
	db := client.Database("eventify").Collection("user")
	var user datatype.User
	err := db.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("Error while FindUser function : %v", err.Error())
	}
	return &user, nil
}

// AddEvents it adds the event to the database
func AddEvents(event datatype.Event) (string, error) {
	db := client.Database("eventify").Collection("event")
	result, err := db.InsertOne(context.TODO(), event)
	if err != nil {
		return "", fmt.Errorf("Error in AddEvents function : %v", err.Error())
	}
	return fmt.Sprintln("Inserted successfully at id : ", result.InsertedID), nil
}

// FindEvents takes user id as parameter finds all the events based on interests if present
func FindEvents(id string) (*[]datatype.Event, error) {
	user, err := FindUser(id)
	if err != nil {
		return nil, err
	}
	db := client.Database("eventify").Collection("event")

	cursor, err := db.Find(context.TODO(), bson.M{"category": bson.M{"$in": user.Interests}})
	var Events []datatype.Event
	ctx := context.TODO()
	if cursor.Next(ctx) {
		var event datatype.Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		Events = append(Events, event)
	}
	if Events == nil {
		event := datatype.Event{}
		Events = append(Events, event)
	}
	return &Events, nil
}

// EditEvent updates an event by its ID
func EditEvent(id string, updatedEvent datatype.Event) error {
	db := client.Database("eventify").Collection("event")
	_, err := db.UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": updatedEvent},
	)
	if err != nil {
		return fmt.Errorf("Error in EditEvent: %v", err.Error())
	}
	return nil
}

// EditUser updates a user by their ID
func EditUser(id string, updatedUser datatype.User) error {
	db := client.Database("eventify").Collection("user")
	_, err := db.UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": updatedUser},
	)
	if err != nil {
		return fmt.Errorf("Error in EditUser: %v", err.Error())
	}
	return nil
}

// AddUserToEvent adds a user ID to the event's participants
func AddUserToEvent(eventID, userID string) error {
	db := client.Database("eventify").Collection("event")
	_, err := db.UpdateOne(
		context.TODO(),
		bson.M{"_id": eventID},
		bson.M{"$addToSet": bson.M{"participants": userID}},
	)
	if err != nil {
		return fmt.Errorf("Error in AddUserToEvent: %v", err.Error())
	}
	return nil
}

// DeleteUser deletes a user by their ID
func DeleteUser(id string) error {
	db := client.Database("eventify").Collection("user")
	_, err := db.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("Error in DeleteUser: %v", err.Error())
	}
	return nil
}

// DeleteEvent deletes an event by its ID
func DeleteEvent(id string) error {
	db := client.Database("eventify").Collection("event")
	_, err := db.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("Error in DeleteEvent: %v", err.Error())
	}
	return nil
}

// RemoveUserFromEvent removes a user ID from the event's participants
func RemoveUserFromEvent(eventID, userID string) error {
	db := client.Database("eventify").Collection("event")
	_, err := db.UpdateOne(
		context.TODO(),
		bson.M{"_id": eventID},
		bson.M{"$pull": bson.M{"participants": userID}},
	)
	if err != nil {
		return fmt.Errorf("Error in RemoveUserFromEvent: %v", err.Error())
	}
	return nil
}
