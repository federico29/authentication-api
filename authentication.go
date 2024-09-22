package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
)

const USERS_TABLE string = "AuthenticationCache"

type AuthUserRequest struct {
	Username string
	Password string
}

type UserResponse struct {
	Username     string
	PasswordHash string
}

func authUser(c *gin.Context) {
	// Get request body
	var request AuthUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user
	user, err := getUser(request.Username)
	if err != nil {
		log.Fatal(err)
		c.String(http.StatusInternalServerError, "Internal server error")
	}
	if user.isEmpty() {
		c.String(http.StatusNotFound, "User Not Found")
		return
	}

	// Auth logic
	requestHash := generateSha256String(request.Password)
	if requestHash == user.PasswordHash {
		c.String(http.StatusOK, "Authorized")
		return
	} else {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
}

func getUser(username string) (UserResponse, error) {
	usernameAttribute, err := attributevalue.Marshal(username)
	if err != nil {
		panic(err)
	}
	key := map[string]types.AttributeValue{"Username": usernameAttribute}
	response, err := dynamoClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: key, TableName: aws.String(USERS_TABLE),
	})
	var userResponse UserResponse
	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", username, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &userResponse)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}

	return userResponse, err
}

func (user UserResponse) isEmpty() bool {
	if user.Username == "" && user.PasswordHash == "" {
		return true
	}
	return false
}

func generateSha256String(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
