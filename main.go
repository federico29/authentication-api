package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var awsRegion string
var dynamoClient *dynamodb.Client

func main() {
	err := loadEnvironment()
	if err != nil {
		log.Fatalf("Error loading environment, %v", err)
		os.Exit(0)
	}

	err = configureAws()
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
		os.Exit(0)
	}

	router := gin.Default()

	router.GET("/", authUser)

	err = router.Run(":5050")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(0)
	}
}

func loadEnvironment() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	awsRegion = os.Getenv("AWS_REGION")

	return nil
}

func configureAws() error {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		return err
	}
	dynamoClient = dynamodb.NewFromConfig(awsConfig)

	return nil
}
