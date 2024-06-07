package routes

import (
	"example.com/goMongo/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/signup", signUp)
	server.POST("/login", logIn)
	server.GET("/getUser", middlewares.Authenticate, getUser)
	server.GET("/getAllUsers", middlewares.Authenticate, getAllUser)
	server.PUT("/updateUser", middlewares.Authenticate, updateUser)
	server.DELETE("/deleteUser", middlewares.Authenticate, deleteUser)

	// Event Routes

	server.POST("/events", middlewares.Authenticate, createEvent)
	server.GET("/events", getEvents)
	server.GET("/events/availableEvents", availableEvents)
	server.GET("/events/:id", getEventByID)
	server.PUT("/events/:id", updateEvent)
	server.DELETE("/events/:id", deleteEvent)
	server.POST("/events/:id/register", middlewares.Authenticate, registerEvent)
	server.GET("/events/registered", middlewares.Authenticate, registeredEvents)
	server.DELETE("events/:id/cancelRegistration", middlewares.Authenticate, cancelRegistration)
}
