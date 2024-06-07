package models

import (
	"context"
	"errors"
	"fmt"

	"example.com/goMongo/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Registration struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EventID primitive.ObjectID `bson:"eventId" json:"eventId"`
	UserID  primitive.ObjectID `bson:"userId" json:"userId"`
	Event   *Event             `bson:"-" json:"event"`
	User    *User              `bson:"-" json:"user"`
}

func RegisterEvent(registration *Registration) (*mongo.InsertOneResult, error) {
	collection := db.GetDatabase().Collection("registrations")
	result, err := collection.InsertOne(context.TODO(), registration)
	if err != nil {
		return nil, err
	}

	// Set the ID field of the user to the inserted ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		registration.ID = oid
	} else {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	return result, nil
}

func RegisteredEvents(userIdStr string) ([]Registration, error) {
	// Convert the userIdStr to primitive.ObjectID
	userIdObj, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		return nil, errors.New("Invalid user ID format")
	}

	// Context to use for the operation
	ctx := context.Background()

	// Get a handle to the collection
	collection := db.GetDatabase().Collection("registrations")

	// Define the filter to only fetch registrations for the logged-in user
	filter := bson.M{"userId": userIdObj}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err // Other error occurred
	}
	defer cursor.Close(ctx)

	var registrations []Registration
	for cursor.Next(ctx) {
		var registration Registration
		if err := cursor.Decode(&registration); err != nil {
			return nil, err // Error decoding event
		}

		// Fetch user data for the event
		userIdHex := registration.UserID.Hex()

		// Fetch user data for the event
		user, err := GetUserById(userIdHex)
		if err != nil {
			return nil, err // Error fetching user data
		}
		registration.User = user

		registrations = append(registrations, registration)
	}

	// Check for any errors that may have occurred during iteration.
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return registrations, nil
}

func CancelRegistration(id string) (*Registration, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("Invalid registration ID format")
	}

	filter := bson.M{"_id": objectID}

	opts := options.FindOneAndDelete()

	ctx := context.Background()

	collection := db.GetDatabase().Collection("registrations")

	var register Registration

	if err := collection.FindOneAndDelete(ctx, filter, opts).Decode(&register); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, err // Other error occurred
	}

	return &register, nil

}
