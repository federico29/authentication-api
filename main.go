package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

var dynamoClient *dynamodb.Client

func main() {
	err := configureAws()
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

func configureAws() error {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return err
	}
	dynamoClient = dynamodb.NewFromConfig(awsConfig)

	return nil
}
