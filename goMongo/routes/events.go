package routes

import (
	"fmt"
	"net/http"

	"example.com/goMongo/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createEvent(c *gin.Context) {

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	fmt.Println(">>>>>>", userId)
	// Assert userId to string
	userIdStr, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User ID format is invalid"})
		return
	}

	// Convert the userIdStr to primitive.ObjectID
	userIdObj, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User ID format is invalid"})
		return
	}

	var event models.Event

	err = c.ShouldBind(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Set the UserID field
	event.UserID = userIdObj
	event.IsAvailable = true

	_, err = models.InsertEvent(&event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "event created successfully", "event": event})
}

func getEvents(c *gin.Context) {

	event, err := models.GetEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to fetch event"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Events fetched", "events": event})
}

func getEventByID(c *gin.Context) {
	eventId := c.Param("id")
	event, err := models.GetEventById(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to fetch event"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, event)
}

func updateEvent(c *gin.Context) {

	eventId := c.Param("id")

	// Parse update data from request body
	var updateData bson.M
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}
	updatedEvent, err := models.UpdateEvent(eventId, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to fetch event"})
		fmt.Println(err)
		return
	}
	if updatedEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "event updated", "event": updatedEvent})
}

func deleteEvent(c *gin.Context) {
	eventId := c.Param("id")

	_, err := models.DeleteEvent(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to fetch user"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "event Deleted"})
}

func availableEvents(c *gin.Context) {
	event, err := models.AvailableEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to fetch event"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Available Events fetched", "events": event})
}
