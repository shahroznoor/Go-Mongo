package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"example.com/goMongo/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `binding:"required" bson:"name" json:"name"`
	Description string             `binding:"required" bson:"description" json:"description"`
	Location    string             `binding:"required" bson:"location" json:"location"`
	DateTime    time.Time          `binding:"required" bson:"dateTime" json:"dateTime"`
	IsAvailable bool               `bson:"isAvailable" json:"isAvailable"`
	UserID      primitive.ObjectID `bson:"userId" json:"userId"` // Reference to the User's ObjectID
	User        *User              `bson:"-" json:"user"`        // Embedded user data
}

func InsertEvent(event *Event) (*mongo.InsertOneResult, error) {
	collection := db.GetDatabase().Collection("events")
	result, err := collection.InsertOne(context.TODO(), event)
	if err != nil {
		return nil, err
	}

	// Set the ID field of the user to the inserted ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		event.ID = oid
	} else {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	return result, nil
}

// Get Events retrieves all events from the MongoDB database.
func GetEvents() ([]Event, error) {
	// Context to use for the operation.
	ctx := context.Background()

	// Get a handle to the collection.
	collection := db.GetDatabase().Collection("events")

	// Perform the query to find all events.
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err // Other error occurred
	}
	defer cursor.Close(ctx)

	var events []Event
	for cursor.Next(ctx) {
		var event Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err // Error decoding event
		}

		// Fetch user data for the event
		userIdHex := event.UserID.Hex()

		// Fetch user data for the event
		user, err := GetUserById(userIdHex)
		if err != nil {
			return nil, err // Error fetching user data
		}
		event.User = user

		events = append(events, event)
	}

	// Check for any errors that may have occurred during iteration.
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func GetEventById(id string) (*Event, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("Invalid user ID format")
	}

	// Specify the filter to find the event by ID.
	filter := bson.M{"_id": objectID}

	// Specify options to configure the query.
	opts := options.FindOne()

	// Context to use for the operation.
	ctx := context.Background()

	// Get a handle to the collection.
	collection := db.GetDatabase().Collection("events")

	// Perform the query.
	var event Event
	if err := collection.FindOne(ctx, filter, opts).Decode(&event); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // event not found
		}
		return nil, err // Other error occurred
	}

	// Fetch user data for the event
	userIdHex := event.UserID.Hex()

	// Fetch user data for the event
	user, err := GetUserById(userIdHex)
	if err != nil {
		return nil, err // Error fetching user data
	}
	event.User = user

	return &event, nil
}

func UpdateEvent(id string, updateData bson.M) (*Event, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("Invalid event ID format")
	}

	filter := bson.M{"_id": objectID}

	// Context to use for the operation.
	ctx := context.Background()

	// Get a handle to the collection.
	collection := db.GetDatabase().Collection("events")

	// Specify the update
	update := bson.M{"$set": updateData}

	// Perform the update.
	var updatedEvent Event
	if err := collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedEvent); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, err // Other error occurred
	}

	return &updatedEvent, nil
}

func DeleteEvent(id string) (*Event, error) {
	// Convert the id string to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("Invalid event ID format")
	}

	// Specify the filter to find the user by ID.
	filter := bson.M{"_id": objectID}

	// Specify options to configure the query.
	opts := options.FindOneAndDelete()

	// Context to use for the operation.
	ctx := context.Background()

	// Get a handle to the collection.
	collection := db.GetDatabase().Collection("events")

	// Perform the query.
	var event Event
	if err := collection.FindOneAndDelete(ctx, filter, opts).Decode(&event); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, err // Other error occurred
	}

	return &event, nil
}

func AvailableEvents() ([]Event, error) {
	// Context to use for the operation.
	ctx := context.Background()

	// Get a handle to the collection.
	collection := db.GetDatabase().Collection("events")

	// Define the filter to only fetch events where isAvailable is true
	filter := bson.M{"isAvailable": true}

	// Perform the query to find all events matching the filter.
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err // Other error occurred
	}
	defer cursor.Close(ctx)

	var events []Event
	for cursor.Next(ctx) {
		var event Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err // Error decoding event
		}

		// Fetch user data for the event
		userIdHex := event.UserID.Hex()

		// Fetch user data for the event
		user, err := GetUserById(userIdHex)
		if err != nil {
			return nil, err // Error fetching user data
		}
		event.User = user

		events = append(events, event)
	}

	// Check for any errors that may have occurred during iteration.
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
