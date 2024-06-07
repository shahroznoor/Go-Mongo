package routes

import (
	"fmt"
	"net/http"

	"example.com/goMongo/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func registerEvent(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

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

	eventId := c.Param("id")
	eventIdObj, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid event ID format"})

	}

	// Check if the event is available
	event, err := models.GetEventById(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve event"})
		return
	}
	if event == nil || !event.IsAvailable {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Event is not available for registration"})
		return
	}

	var registration models.Registration
	registration.EventID = eventIdObj
	registration.UserID = userIdObj

	_, err = models.RegisterEvent(&registration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Update the event to set IsAvailable to false
	updateData := bson.M{"isAvailable": false}
	_, err = models.UpdateEvent(eventId, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update event availability"})
		return
	}

	registration.Event = event
	user, err := models.GetUserById(userIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user"})
		return
	}
	registration.User = user

	c.JSON(http.StatusCreated, gin.H{"message": "event Registered successfully", "registration": registration})

}

func registeredEvents(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID not found in context"})
		return
	}

	// Assert userId to string
	userIdStr, ok := userId.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User ID format is invalid"})
		return
	}

	events, err := models.RegisteredEvents(userIdStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"registrations": events})
}

func cancelRegistration(c *gin.Context) {
	registrationId := c.Param("id")

	_, err := models.CancelRegistration(registrationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to fetch user"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "registration cancelled"})
}
