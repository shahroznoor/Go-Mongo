package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const secretKey = "secretkey"

func GenerateToken(email string, userId primitive.ObjectID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId.Hex(), // Store ObjectID as a string
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(token string) (primitive.ObjectID, error) {
	// Parse the JWT token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is HMAC
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Return the secret key for validation
		return []byte(secretKey), nil
	})
	if err != nil {
		return primitive.NilObjectID, errors.New("could not parse token")
	}

	// Check if the token is valid
	if !parsedToken.Valid {
		return primitive.NilObjectID, errors.New("invalid token")
	}

	// Extract claims from the token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return primitive.NilObjectID, errors.New("invalid token claims")
	}

	// Extract the user ID from claims
	userIdHex, ok := claims["userId"].(string)
	if !ok {
		return primitive.NilObjectID, errors.New("invalid user ID in token claims")
	}

	// Convert user ID string to ObjectID
	userId, err := primitive.ObjectIDFromHex(userIdHex)
	if err != nil {
		return primitive.NilObjectID, errors.New("invalid user ID format")
	}

	return userId, nil
}
