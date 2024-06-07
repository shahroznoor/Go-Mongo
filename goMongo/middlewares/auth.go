package middlewares

import (
	"net/http"

	"example.com/goMongo/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	// Get the token from the request header
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token is required"})
		return
	}

	// Verify the token and extract the user ID
	userId, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Set the user ID in the Gin context
	context.Set("userId", userId.Hex()) // Convert ObjectID to string

	context.Next()
}
