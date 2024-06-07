package models

import (
	"context"
	"errors"
	"fmt"

	"example.com/goMongo/db"
	"example.com/goMongo/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User represents a user in the system
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `binding:required bson:"email" json:"email"`
	Password string             `binding:required bson:"password" json:"password"`
}

// InsertUser inserts a new user into the database
func InsertUser(user *User) (*mongo.InsertOneResult, error) {

	// Check if the email already exists
	if emailExists(user.Email) {
		return nil, errors.New("Email already exists")
	}
	// Hash the user's password before inserting
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	collection := db.GetDatabase().Collection("users")
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	// Set the ID field of the user to the inserted ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	} else {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	return result, nil
}

// emailExists checks if the given email already exists in the database
func emailExists(email string) bool {
	collection := db.GetDatabase().Collection("users")

	filter := bson.M{"email": email}

	var existingUser User
	err := collection.FindOne(context.TODO(), filter).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		return false // Email does not exist
	} else if err != nil {
		fmt.Println("Error checking email existence:", err)
		return true // Assume email exists to avoid duplicate inserts
	}

	return true // Email exists
}
func (u *User) ValidateCredentials() error {
	// Search for the user by email
	filter := bson.M{"email": u.Email}
	collection := db.GetDatabase().Collection("users")
	var userFromDB User
	err := collection.FindOne(context.TODO(), filter).Decode(&userFromDB)
	if err == mongo.ErrNoDocuments {
		return errors.New("Invalid email")
	} else if err != nil {
		return err
	}

	// Compare the provided password with the retrieved password
	passwordIsValid := utils.CheckPassword(u.Password, userFromDB.Password)
	if !passwordIsValid {
		return errors.New("Invalid password")
	}

	// Set the user ID from the retrieved user
	u.ID = userFromDB.ID

	return nil
}

// GetUserById retrieves a user from the MongoDB database by ID.
func GetUserById(id string) (*User, error) {
	// Convert the id string to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("Invalid user ID format")
	}

	// Specify the filter to find the user by ID.
	filter := bson.M{"_id": objectID}

	// Specify options to configure the query.
	opts := options.FindOne()

	// Context to use for the operation.
	ctx := context.Background()

	// Get a handle to the collection.
	collection := db.GetDatabase().Collection("users")

	// Perform the query.
	var user User
	if err := collection.FindOne(ctx, filter, opts).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, err // Other error occurred
	}

	return &user, nil
}

// GetUsers retrieves all users from the MongoDB database.
func GetUsers() ([]User, error) {
	// Context to use for the operation.
	ctx := context.Background()

	// Get a handle to the collection.
	collection := db.GetDatabase().Collection("users")

	// Perform the query to find all users.
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err // Other error occurred
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err // Error decoding user
		}
		users = append(users, user)
	}

	// Check for any errors that may have occurred during iteration.
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUserById updates a user's details in the MongoDB database by ID.
func UpdateUserById(id string, updateData bson.M) (*User, error) {
	// Convert the id string to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("Invalid user ID format")
	}

	// Check if the email already exists if the email is being updated
	if newEmail, ok := updateData["email"].(string); ok && emailExists(newEmail) {
		return nil, errors.New("Email already exists")
	}

	// Specify the filter to find the user by ID.
	filter := bson.M{"_id": objectID}

	// Context to use for the operation.
	ctx := context.Background()

	// Get a handle to the collection.
	collection := db.GetDatabase().Collection("users")

	// Specify the update
	update := bson.M{"$set": updateData}

	// Perform the update.
	var updatedUser User
	if err := collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedUser); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, err // Other error occurred
	}

	return &updatedUser, nil
}

// GetUserById retrieves a user from the MongoDB database by ID.
func DeleteUser(id string) (*User, error) {
	// Convert the id string to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("Invalid user ID format")
	}

	// Specify the filter to find the user by ID.
	filter := bson.M{"_id": objectID}

	// Specify options to configure the query.
	opts := options.FindOneAndDelete()

	// Context to use for the operation.
	ctx := context.Background()

	// Get a handle to the collection.
	collection := db.GetDatabase().Collection("users")

	// Perform the query.
	var user User
	if err := collection.FindOneAndDelete(ctx, filter, opts).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, err // Other error occurred
	}

	return &user, nil
}
