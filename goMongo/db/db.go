package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var database *mongo.Database

// InitDB initializes the database connection
func InitDB() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	// Get a handle for your database
	database = client.Database("api_db")
	fmt.Println("Connected to MongoDB!")
}

// GetDatabase returns the database handle
func GetDatabase() *mongo.Database {
	if database == nil {
		InitDB()
	}
	return database
}
