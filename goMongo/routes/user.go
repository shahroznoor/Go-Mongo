package routes

import (
	"fmt"
	"net/http"

	"example.com/goMongo/models"
	"example.com/goMongo/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func signUp(c *gin.Context) {
	var user models.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	result, err := models.InsertUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	fmt.Println("Inserted user with ID:", result.InsertedID)

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		fmt.Println(">>>>", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "JWT token generating error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "signup successfully", "user": user, "token": token})
}

func logIn(c *gin.Context) {
	var user models.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	err = user.ValidateCredentials()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		fmt.Println(">>>>", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "JWT token generating error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user logged In", "token": token})
}

func getAllUser(c *gin.Context) {

	user, err := models.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to fetch user"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Users fetched", "users": user})
}

func getUser(c *gin.Context) {
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

	user, err := models.GetUserById(userIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to fetch user"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User fetched", "user": user})
}

func updateUser(c *gin.Context) {
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

	// Parse update data from request body
	var updateData bson.M
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	// Call the UpdateUserById function
	updatedUser, err := models.UpdateUserById(userIdStr, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		fmt.Println(err)
		return
	}
	if updatedUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func deleteUser(c *gin.Context) {
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

	_, err := models.DeleteUser(userIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to fetch user"})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user Deleted"})
}
